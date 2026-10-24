// Package eval runs the retrieval/answer experiments for the research component (R1–R9):
// a test set, runs per retrieval mode, Recall@K and MRR, expert grading and CSV export.
package eval

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reclover/sinakhontao/internal/httpx"
	"github.com/reclover/sinakhontao/internal/rag"
)

// RecallAtK is the share of relevant documents found in the top k retrieved.
func RecallAtK(retrieved, relevant []string, k int) float64 {
	if len(relevant) == 0 {
		return 0
	}
	top := map[string]bool{}
	for i, d := range retrieved {
		if i >= k {
			break
		}
		top[d] = true
	}
	hit := 0
	for _, d := range relevant {
		if top[d] {
			hit++
		}
	}
	return float64(hit) / float64(len(relevant))
}

// ReciprocalRank is 1/rank of the first relevant document, or 0.
func ReciprocalRank(retrieved, relevant []string) float64 {
	rel := map[string]bool{}
	for _, d := range relevant {
		rel[d] = true
	}
	for i, d := range retrieved {
		if rel[d] {
			return 1 / float64(i+1)
		}
	}
	return 0
}

type Handler struct {
	DB      *pgxpool.Pool
	RAG     *rag.Service
	running atomic.Bool // one experiment at a time bounds API cost (threat A4)
}

func (h *Handler) ListQuestions(w http.ResponseWriter, r *http.Request) error {
	return rawJSON(w, r, h.DB, `SELECT coalesce(jsonb_agg(to_jsonb(q) ORDER BY q.id), '[]') FROM eval_questions q`)
}

// ImportQuestions accepts the AI_Test_Set sheet as CSV:
// test_id, lang, question, reference_answer, relevant_doc_ids (separated by ;), category, difficulty, reviewer
func (h *Handler) ImportQuestions(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	cr := csv.NewReader(r.Body)
	cr.FieldsPerRecord = -1
	header, err := cr.Read()
	if err != nil {
		return httpx.BadRequest("empty CSV")
	}
	idx := map[string]int{}
	for i, c := range header {
		idx[strings.TrimSpace(strings.TrimPrefix(c, "\xEF\xBB\xBF"))] = i
	}
	for _, c := range []string{"test_id", "lang", "question", "reference_answer", "relevant_doc_ids"} {
		if _, ok := idx[c]; !ok {
			return httpx.BadRequest("missing column " + c)
		}
	}
	get := func(rec []string, col string) string {
		if i, ok := idx[col]; ok && i < len(rec) {
			return strings.TrimSpace(rec[i])
		}
		return ""
	}
	tx, err := h.DB.Begin(r.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(r.Context())
	n := 0
	for line := 2; ; line++ {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return httpx.BadRequest(fmt.Sprintf("line %d: %v", line, err))
		}
		lang := get(rec, "lang")
		if lang != "th" && lang != "zh" && lang != "en" {
			return httpx.BadRequest(fmt.Sprintf("line %d: lang must be th, zh or en", line))
		}
		var docs []string
		for _, d := range strings.Split(get(rec, "relevant_doc_ids"), ";") {
			if d = strings.TrimSpace(d); d != "" {
				docs = append(docs, d)
			}
		}
		if get(rec, "test_id") == "" || get(rec, "question") == "" || len(docs) == 0 {
			return httpx.BadRequest(fmt.Sprintf("line %d: test_id, question and relevant_doc_ids are required", line))
		}
		diff := get(rec, "difficulty")
		if diff == "" {
			diff = "medium"
		}
		_, err = tx.Exec(r.Context(), `INSERT INTO eval_questions(code, lang, question, reference_answer, relevant_docs, category, difficulty, reviewer)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (code) DO UPDATE SET lang=EXCLUDED.lang, question=EXCLUDED.question,
			reference_answer=EXCLUDED.reference_answer, relevant_docs=EXCLUDED.relevant_docs, category=EXCLUDED.category,
			difficulty=EXCLUDED.difficulty, reviewer=EXCLUDED.reviewer`,
			get(rec, "test_id"), lang, get(rec, "question"), get(rec, "reference_answer"), docs, get(rec, "category"), diff, get(rec, "reviewer"))
		if err != nil {
			return err
		}
		n++
	}
	if err := tx.Commit(r.Context()); err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusOK, map[string]int{"imported": n})
}

func (h *Handler) ListRuns(w http.ResponseWriter, r *http.Request) error {
	return rawJSON(w, r, h.DB, `SELECT coalesce(jsonb_agg(to_jsonb(x) ORDER BY x.id DESC), '[]') FROM eval_runs x`)
}

func (h *Handler) StartRun(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Mode        string `json:"mode"`
		K           int    `json:"k"`
		Note        string `json:"note"`
		WithAnswers bool   `json:"with_answers"` // generating answers calls the LLM: costs money
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	if !rag.ValidMode(in.Mode) {
		return httpx.BadRequest("mode must be llm_only, vector, multilingual or hybrid")
	}
	if in.K < 1 || in.K > 20 {
		return httpx.BadRequest("k must be 1-20")
	}
	if len(in.Note) > 500 {
		return httpx.BadRequest("note too long")
	}
	if !h.running.CompareAndSwap(false, true) {
		return &httpx.Error{Status: http.StatusConflict, Code: "busy", Message: "another evaluation run is in progress"}
	}
	var id int64
	if err := h.DB.QueryRow(r.Context(), `INSERT INTO eval_runs(mode, k, note) VALUES ($1,$2,$3) RETURNING id`, in.Mode, in.K, in.Note).Scan(&id); err != nil {
		h.running.Store(false)
		return err
	}
	go h.run(id, rag.Mode(in.Mode), in.K, in.WithAnswers)
	return httpx.JSON(w, http.StatusAccepted, map[string]int64{"id": id})
}

func (h *Handler) run(runID int64, mode rag.Mode, k int, withAnswers bool) {
	defer h.running.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()
	fail := func(err error) {
		slog.Error("eval run failed", "run", runID, "err", err)
		h.DB.Exec(ctx, `UPDATE eval_runs SET status='failed', note = note || ' | ' || $2 WHERE id=$1`, runID, err.Error())
	}
	rows, err := h.DB.Query(ctx, `SELECT id, lang, question, relevant_docs FROM eval_questions ORDER BY id`)
	if err != nil {
		fail(err)
		return
	}
	type q struct {
		id       int64
		lang, qs string
		rel      []string
	}
	var qs []q
	for rows.Next() {
		var x q
		if err := rows.Scan(&x.id, &x.lang, &x.qs, &x.rel); err != nil {
			rows.Close()
			fail(err)
			return
		}
		qs = append(qs, x)
	}
	rows.Close()

	var sumRecall, sumRR float64
	for _, x := range qs {
		start := time.Now()
		var docs []string
		var answer string
		if withAnswers || mode == rag.ModeNone {
			a, err := h.RAG.Ask(ctx, mode, x.lang, x.qs, k)
			if err != nil {
				fail(err)
				return
			}
			docs, answer = rag.DocRanking(a.Retrieved), a.Answer
		} else {
			hits, err := h.RAG.Retriever.Retrieve(ctx, mode, x.lang, x.qs, k)
			if err != nil {
				fail(err)
				return
			}
			docs = rag.DocRanking(hits)
		}
		if docs == nil {
			docs = []string{}
		}
		rec, rr := RecallAtK(docs, x.rel, k), ReciprocalRank(docs, x.rel)
		sumRecall += rec
		sumRR += rr
		if _, err := h.DB.Exec(ctx, `INSERT INTO eval_results(run_id, question_id, retrieved_docs, recall, rr, answer, latency_ms) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			runID, x.id, docs, rec, rr, answer, int(time.Since(start).Milliseconds())); err != nil {
			fail(err)
			return
		}
	}
	n := float64(max(len(qs), 1))
	h.DB.Exec(ctx, `UPDATE eval_runs SET status='done', n=$2, recall_at_k=$3, mrr=$4 WHERE id=$1`, runID, len(qs), sumRecall/n, sumRR/n)
}

func (h *Handler) RunResults(w http.ResponseWriter, r *http.Request) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	var out json.RawMessage
	err = h.DB.QueryRow(r.Context(), `SELECT coalesce(jsonb_agg(jsonb_build_object(
		'id', e.id, 'code', q.code, 'lang', q.lang, 'question', q.question, 'reference_answer', q.reference_answer,
		'relevant_docs', q.relevant_docs, 'retrieved_docs', e.retrieved_docs, 'recall', e.recall, 'rr', e.rr,
		'answer', e.answer, 'latency_ms', e.latency_ms, 'expert_correct', e.expert_correct,
		'expert_faithful', e.expert_faithful, 'expert_note', e.expert_note) ORDER BY e.id), '[]')
		FROM eval_results e JOIN eval_questions q ON q.id = e.question_id WHERE e.run_id = $1`, id).Scan(&out)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(out)
	return err
}

// Grade stores an expert's 1–5 correctness and faithfulness scores for one answer (R6, R7).
func (h *Handler) Grade(w http.ResponseWriter, r *http.Request) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	var in struct {
		Correct  int    `json:"correct"`
		Faithful int    `json:"faithful"`
		Note     string `json:"note"`
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	if in.Correct < 1 || in.Correct > 5 || in.Faithful < 1 || in.Faithful > 5 || len(in.Note) > 2000 {
		return httpx.BadRequest("scores must be 1-5")
	}
	tag, err := h.DB.Exec(r.Context(), `UPDATE eval_results SET expert_correct=$2, expert_faithful=$3, expert_note=$4 WHERE id=$1`, id, in.Correct, in.Faithful, in.Note)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httpx.NotFound()
	}
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) ExportRuns(w http.ResponseWriter, r *http.Request) error {
	rows, err := h.DB.Query(r.Context(), `SELECT r.id, r.mode, r.k, q.code, q.lang, q.category, q.difficulty,
		array_to_string(q.relevant_docs, ';'), array_to_string(e.retrieved_docs, ';'), e.recall, e.rr, e.latency_ms,
		coalesce(e.expert_correct::text,''), coalesce(e.expert_faithful::text,''), e.answer
		FROM eval_results e JOIN eval_runs r ON r.id = e.run_id JOIN eval_questions q ON q.id = e.question_id
		ORDER BY r.id, q.id`)
	if err != nil {
		return err
	}
	defer rows.Close()
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="eval-results.csv"`)
	w.Write([]byte("\xEF\xBB\xBF"))
	cw := csv.NewWriter(w)
	cw.Write([]string{"run_id", "mode", "k", "test_id", "lang", "category", "difficulty", "relevant_docs", "retrieved_docs", "recall_at_k", "reciprocal_rank", "latency_ms", "expert_correct", "expert_faithful", "answer"})
	for rows.Next() {
		var runID int64
		var k, lat int
		var mode, code, lang, cat, diff, rel, ret, ec, ef, ans string
		var rec, rr float64
		if err := rows.Scan(&runID, &mode, &k, &code, &lang, &cat, &diff, &rel, &ret, &rec, &rr, &lat, &ec, &ef, &ans); err != nil {
			return err
		}
		cw.Write([]string{strconv.FormatInt(runID, 10), mode, strconv.Itoa(k), httpx.CSVSafe(code), lang, httpx.CSVSafe(cat), diff, httpx.CSVSafe(rel), httpx.CSVSafe(ret),
			strconv.FormatFloat(rec, 'f', 4, 64), strconv.FormatFloat(rr, 'f', 4, 64), strconv.Itoa(lat), ec, ef, httpx.CSVSafe(ans)})
	}
	cw.Flush()
	return rows.Err()
}

func rawJSON(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool, q string) error {
	var out json.RawMessage
	if err := db.QueryRow(r.Context(), q).Scan(&out); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(out)
	return err
}
