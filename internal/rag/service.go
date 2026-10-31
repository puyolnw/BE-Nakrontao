package rag

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reclover/sinakhontao/internal/httpx"
)

// Evidence gate (scope 5.10): answer only when a chunk is semantically close enough OR shares a
// strong keyword match. The vector threshold depends on the embedder, so it is configurable.
//
// Measured on the demo corpus (th/zh/en): trigram word similarity is <= 0.13 for off-topic
// questions and >= 0.40 for on-topic ones in every language, so it is the primary signal.
// The offline hash embedder cannot separate Thai topics reliably, so with it a keyword match is
// mandatory; a real embedding model may also pass the gate on semantic similarity alone.
const (
	MinKeywordSim          = 0.30
	DefaultMinVecSimRemote = 0.30
)

type Service struct {
	DB          *pgxpool.Pool
	Retriever   *Retriever
	LLM         LLM
	Salt        string
	MinVecSim   float64 // 0 disables the semantic-only path (offline hash embedder)
	PerClient   *httpx.Limiter
	Daily       *httpx.DailyCap      // whole site, per UTC day
	ClientDaily *httpx.KeyedDailyCap // one visitor, per UTC day

	indexMu   sync.Mutex
	slots     chan struct{}
	slotsOnce sync.Once
}

// MaxConcurrentLLM bounds simultaneous LLM calls so a burst cannot run up cost or exhaust memory (threat A4).
const MaxConcurrentLLM = 3

var ErrBusy = &httpx.Error{Status: http.StatusServiceUnavailable, Code: "busy", Message: "assistant is busy, please retry shortly"}

func (s *Service) acquire() (func(), bool) {
	s.slotsOnce.Do(func() { s.slots = make(chan struct{}, MaxConcurrentLLM) })
	select {
	case s.slots <- struct{}{}:
		return func() { <-s.slots }, true
	default:
		return nil, false
	}
}

type Answer struct {
	ID        int64  `json:"id,omitempty"`
	Answer    string `json:"answer"`
	Citations []Hit  `json:"citations"`
	Retrieved []Hit  `json:"-"` // full ranked retrieval, used for research metrics
	Refused   bool   `json:"refused"`
	LatencyMS int    `json:"latency_ms"`
}

var refusals = map[string]string{
	"th": "ขออภัย ยังไม่มีข้อมูลที่ได้รับการตรวจสอบเพียงพอสำหรับคำถามนี้ กรุณาสอบถามเจ้าหน้าที่ศาลเจ้าหรือผู้นำชุมชน",
	"zh": "抱歉，目前没有足够的经过核实的资料回答这个问题。请咨询神祠工作人员或社区负责人。",
	"en": "Sorry, there is not enough verified information to answer this. Please ask the shrine staff or community leaders.",
}

// Ask retrieves evidence and generates an answer. It is shared by the chatbot and the eval runner.
func (s *Service) Ask(ctx context.Context, mode Mode, lang, question string, k int) (Answer, error) {
	return s.ask(ctx, systemPrompt, mode, lang, question, k)
}

func (s *Service) ask(ctx context.Context, system string, mode Mode, lang, question string, k int) (Answer, error) {
	start := time.Now()
	hits, err := s.Retriever.Retrieve(ctx, mode, lang, question, k)
	if err != nil {
		return Answer{}, err
	}
	a := Answer{Citations: hits, Retrieved: hits}
	if a.Citations == nil {
		a.Citations = []Hit{}
	}
	if mode != ModeNone && !strongEnough(mode, hits, s.MinVecSim) {
		a.Answer, a.Refused = refusals[lang], true
		a.LatencyMS = int(time.Since(start).Milliseconds())
		return a, nil
	}
	var text string
	release, ok := s.acquire()
	if !ok {
		return Answer{}, ErrBusy
	}
	defer release()
	if mode == ModeNone {
		text, err = s.LLM.Complete(ctx, system+"\n- No passages are given in this condition: answer from general knowledge, and say NO_ANSWER if unsure.",
			"<lang>"+langName(lang)+"</lang>\n<question>"+escapeForPrompt(question)+"</question>")
	} else {
		text, err = s.LLM.Complete(ctx, system, buildPrompt(lang, question, hits))
	}
	if err != nil {
		return Answer{}, err
	}
	if text == "" || strings.Contains(text, "NO_ANSWER") {
		a.Answer, a.Refused = refusals[lang], true
	} else {
		a.Answer = normalizeCitations(text)
		a.Citations = citedHits(a.Answer, hits, s.LLM.Name() == "extractive")
	}
	a.LatencyMS = int(time.Since(start).Milliseconds())
	return a, nil
}

var citationMark = regexp.MustCompile(`[\[【]\s*([A-Za-z][A-Za-z0-9-]*)\s*[\]】]`)

// normalizeCitations rewrites the variants models emit, such as 【K-001】 or [ K-001 ], to [K-001].
func normalizeCitations(s string) string { return citationMark.ReplaceAllString(s, "[$1]") }

// citedHits keeps only the passages the answer actually relies on: those whose id appears as
// [DOC-CODE] in the generated text, or the quoted first passage for the extractive fallback.
func citedHits(answer string, hits []Hit, extractive bool) []Hit {
	if extractive && len(hits) > 0 {
		return hits[:1]
	}
	var out []Hit
	seen := map[string]bool{}
	for _, h := range hits {
		if !seen[h.DocCode] && strings.Contains(answer, "["+h.DocCode+"]") {
			seen[h.DocCode] = true
			out = append(out, h)
		}
	}
	if out == nil {
		return []Hit{}
	}
	return out
}

func strongEnough(mode Mode, hits []Hit, minVec float64) bool {
	for _, h := range hits {
		if h.KwSim >= MinKeywordSim {
			return true
		}
		// Vector-only research modes have no keyword signal; fall back to the vector threshold.
		vecFloor := minVec
		if vecFloor == 0 && mode != ModeHybrid {
			vecFloor = 0.25
		}
		if vecFloor > 0 && h.VecSim >= vecFloor {
			return true
		}
	}
	return false
}

func (s *Service) Chat(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		SessionID string `json:"session_id"`
		Lang      string `json:"lang"`
		Question  string `json:"question"`
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	in.Question = strings.TrimSpace(in.Question)
	if n := utf8.RuneCountInString(in.Question); n < 2 || n > 500 {
		return httpx.BadRequest("question must be 2-500 characters")
	}
	if in.Lang != "zh" && in.Lang != "en" {
		in.Lang = "th"
	}
	if len(in.SessionID) < 8 || len(in.SessionID) > 64 {
		return httpx.BadRequest("invalid session_id")
	}
	visitor := httpx.HashIP(s.Salt, httpx.ClientIP(r))
	if !s.PerClient.Allow(visitor) {
		return httpx.TooMany()
	}
	// PII never leaves for the LLM provider, and only the redacted text is logged.
	question := RedactPII(in.Question)
	in.Lang = detectLang(question, in.Lang)
	var a Answer
	if reply, ok := screen(in.Lang, question); ok {
		a = Answer{Answer: reply, Refused: true}
	} else {
		if !s.ClientDaily.Allow(visitor) {
			return &httpx.Error{Status: http.StatusTooManyRequests, Code: "daily_limit", Message: "daily question limit reached, please try tomorrow"}
		}
		if !s.Daily.Take() {
			return &httpx.Error{Status: http.StatusServiceUnavailable, Code: "quota", Message: "daily AI quota reached, please try tomorrow"}
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		var err error
		if a, err = s.chatAnswer(ctx, in.Lang, question); err != nil {
			return err
		}
	}
	if a.Citations == nil {
		a.Citations = []Hit{} // a refusal has none; the column is NOT NULL
	}
	cites, _ := json.Marshal(a.Citations)
	err := s.DB.QueryRow(r.Context(), `INSERT INTO chat_logs(session_id, lang, question, answer, citations, refused, category, latency_ms)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		in.SessionID, in.Lang, question, RedactPII(a.Answer), cites, a.Refused, Classify(question), a.LatencyMS).Scan(&a.ID)
	if err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusOK, a)
}

func (s *Service) Feedback(w http.ResponseWriter, r *http.Request) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	var in struct {
		SessionID string `json:"session_id"`
		Helpful   bool   `json:"helpful"`
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	tag, err := s.DB.Exec(r.Context(), `UPDATE chat_logs SET helpful=$3 WHERE id=$1 AND session_id=$2`, id, in.SessionID, in.Helpful)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httpx.NotFound()
	}
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Reindex rebuilds kb_chunks from published knowledge and FAQ rows (scope 5.5, 5.11).
// The knowledge base is small (hundreds of docs), so a full rebuild is simpler and safer than diffing.
func (s *Service) Reindex(ctx context.Context) error {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()
	type doc struct{ code, lang, kind, title, body, source string }
	var docs []doc
	rows, err := s.DB.Query(ctx, `
		SELECT k.doc_code, l.key, k.kind, coalesce(k.title->>l.key, k.title->>'th'), l.value, k.source
		FROM knowledge k, jsonb_each_text(k.body) l
		WHERE k.status = 'published' AND l.value <> ''
		UNION ALL
		SELECT coalesce(nullif(f.doc_code,''), 'FAQ-' || f.id), l.key, 'faq', f.question->>l.key, (f.question->>l.key) || E'\n' || l.value, 'FAQ'
		FROM faq f, jsonb_each_text(f.answer) l
		WHERE f.published AND l.value <> '' AND coalesce(f.question->>l.key, '') <> ''`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var d doc
		if err := rows.Scan(&d.code, &d.lang, &d.kind, &d.title, &d.body, &d.source); err != nil {
			rows.Close()
			return err
		}
		docs = append(docs, d)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	type chunk struct {
		doc
		content string
	}
	var chunks []chunk
	for _, d := range docs {
		for _, c := range Chunk(d.body) {
			chunks = append(chunks, chunk{d, c})
		}
	}
	texts := make([]string, len(chunks))
	for i, c := range chunks {
		texts[i] = c.title + "\n" + c.content
	}
	var vecs [][]float32
	for i := 0; i < len(texts); i += 64 {
		end := min(i+64, len(texts))
		v, err := s.Retriever.Embedder.Embed(ctx, texts[i:end])
		if err != nil {
			return err
		}
		vecs = append(vecs, v...)
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `DELETE FROM kb_chunks`); err != nil {
		return err
	}
	for i, c := range chunks {
		if _, err := tx.Exec(ctx, `INSERT INTO kb_chunks(doc_code, lang, kind, title, content, source, embedding) VALUES ($1,$2,$3,$4,$5,$6,$7::vector)`,
			c.code, c.lang, c.kind, c.title, c.content, c.source, vectorLiteral(vecs[i])); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	slog.Info("knowledge base reindexed", "docs", len(docs), "chunks", len(chunks), "embedder", s.Retriever.Embedder.Name())
	return nil
}

// ReindexAsync is the hook used after CMS edits; it detaches from the request context.
func (s *Service) ReindexAsync(_ context.Context) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := s.Reindex(ctx); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("reindex failed", "err", err)
		}
	}()
}

func (s *Service) ReindexHandler(w http.ResponseWriter, r *http.Request) error {
	if err := s.Reindex(r.Context()); err != nil {
		return err
	}
	var n int
	s.DB.QueryRow(r.Context(), `SELECT count(*) FROM kb_chunks`).Scan(&n)
	return httpx.JSON(w, http.StatusOK, map[string]int{"chunks": n})
}

// Purge deletes chat logs and page views older than the retention period (threat A5, PDPA).
func Purge(ctx context.Context, db *pgxpool.Pool, days int) {
	db.Exec(ctx, `DELETE FROM chat_logs WHERE created_at < now() - make_interval(days => $1)`, days)
	db.Exec(ctx, `DELETE FROM page_views WHERE created_at < now() - make_interval(days => $1)`, days)
	db.Exec(ctx, `DELETE FROM sessions WHERE expires_at < now()`)
}

// Classify assigns a coarse topic to a question for the popular-categories report (scope 9.6).
func Classify(q string) string {
	q = strings.ToLower(q)
	topics := []struct {
		name  string
		words []string
	}{
		{"ritual", []string{"ไหว้", "สักการะ", "บูชา", "ของไหว้", "ขอพร", "แก้บน", "พิธี", "祭拜", "拜", "供品", "祈福", "仪式", "worship", "offering", "pray", "ritual"}},
		{"history", []string{"ประวัติ", "ตำนาน", "ความเป็นมา", "โบราณ", "历史", "传说", "来历", "history", "legend", "origin"}},
		{"event", []string{"งาน", "เทศกาล", "วันไหน", "ประเพณี", "节日", "活动", "庆典", "festival", "event", "when"}},
		{"travel", []string{"เดินทาง", "ไปยังไง", "ที่จอดรถ", "เปิด", "ปิด", "เส้นทาง", "交通", "怎么去", "开放", "路线", "how to get", "get there", "how do i get", "open", "route", "parking"}},
		{"community", []string{"ร้าน", "อาหาร", "ที่พัก", "สินค้า", "ของฝาก", "餐", "住宿", "特产", "商店", "food", "shop", "hotel", "souvenir"}},
	}
	for _, t := range topics {
		for _, w := range t.words {
			if strings.Contains(q, w) {
				return t.name
			}
		}
	}
	return "other"
}
