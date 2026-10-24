package insights

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reclover/sinakhontao/internal/auth"
	"github.com/reclover/sinakhontao/internal/httpx"
)

type Handler struct {
	DB    *pgxpool.Pool
	Salt  string
	Track *httpx.Limiter
	Votes *httpx.Limiter
}

// ---- public ----

func (h *Handler) SurveyQuestions(w http.ResponseWriter, r *http.Request) error {
	var out json.RawMessage
	err := h.DB.QueryRow(r.Context(), `SELECT coalesce(jsonb_agg(jsonb_build_object('id',id,'section',section,'text',text) ORDER BY section, sort, id), '[]')
		FROM survey_questions WHERE active`).Scan(&out)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(out)
	return err
}

type surveyInput struct {
	Lang       string            `json:"lang"`
	Consent    bool              `json:"consent"`
	Respondent map[string]string `json:"respondent"` // coarse demographics only: age_range, nationality, visitor_type
	Answers    map[string]int    `json:"answers"`    // question id -> 1..5
	Comment    string            `json:"comment"`
}

var respondentKeys = map[string]bool{"age_range": true, "gender": true, "nationality": true, "visitor_type": true}

func (h *Handler) SubmitSurvey(w http.ResponseWriter, r *http.Request) error {
	var in surveyInput
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	if !in.Consent {
		return httpx.BadRequest("consent is required before collecting research data")
	}
	if !h.Votes.Allow(httpx.HashIP(h.Salt, httpx.ClientIP(r))) {
		return httpx.TooMany()
	}
	if in.Lang != "zh" && in.Lang != "en" {
		in.Lang = "th"
	}
	if len(in.Answers) == 0 || len(in.Answers) > 100 {
		return httpx.BadRequest("answers are required")
	}
	for k, v := range in.Answers {
		if _, err := strconv.ParseInt(k, 10, 64); err != nil || v < 1 || v > 5 {
			return httpx.BadRequest("answers must map question ids to 1-5")
		}
	}
	for k, v := range in.Respondent {
		if !respondentKeys[k] || len(v) > 60 {
			return httpx.BadRequest("unsupported respondent field " + k)
		}
	}
	if len(in.Comment) > 4000 {
		return httpx.BadRequest("comment too long")
	}
	resp, _ := json.Marshal(in.Respondent)
	ans, _ := json.Marshal(in.Answers)
	if _, err := h.DB.Exec(r.Context(), `INSERT INTO survey_responses(lang, consent, respondent, answers, comment) VALUES ($1,true,$2,$3,$4)`,
		in.Lang, resp, ans, strings.TrimSpace(in.Comment)); err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h *Handler) TrackView(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Path string `json:"path"`
		Lang string `json:"lang"`
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	visitor := httpx.HashIP(h.Salt, httpx.ClientIP(r))
	if !h.Track.Allow(visitor) || len(in.Path) == 0 || len(in.Path) > 200 || !strings.HasPrefix(in.Path, "/") {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	if in.Lang != "zh" && in.Lang != "en" {
		in.Lang = "th"
	}
	if _, err := h.DB.Exec(r.Context(), `INSERT INTO page_views(path, lang, visitor) VALUES ($1,$2,$3)`, in.Path, in.Lang, visitor); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// ---- admin ----

func dateRange(r *http.Request) (time.Time, time.Time) {
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	if t, err := time.Parse("2006-01-02", r.URL.Query().Get("from")); err == nil {
		from = t
	}
	if t, err := time.Parse("2006-01-02", r.URL.Query().Get("to")); err == nil {
		to = t.AddDate(0, 0, 1)
	}
	return from, to
}

type kv struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

func (h *Handler) groupCount(r *http.Request, q string, args ...any) ([]kv, error) {
	rows, err := h.DB.Query(r.Context(), q, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (kv, error) {
		var x kv
		return x, row.Scan(&x.Key, &x.Count)
	})
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) error {
	from, to := dateRange(r)
	out := map[string]any{"from": from.Format("2006-01-02"), "to": to.AddDate(0, 0, -1).Format("2006-01-02")}

	var totals struct {
		Places, Knowledge, Events, News, Shops, Chunks int
	}
	err := h.DB.QueryRow(r.Context(), `SELECT
		(SELECT count(*) FROM places), (SELECT count(*) FROM knowledge), (SELECT count(*) FROM events),
		(SELECT count(*) FROM news), (SELECT count(*) FROM shops), (SELECT count(*) FROM kb_chunks)`).
		Scan(&totals.Places, &totals.Knowledge, &totals.Events, &totals.News, &totals.Shops, &totals.Chunks)
	if err != nil {
		return err
	}
	out["totals"] = totals

	var views, visitors, chats, refused, helpful, unhelpful, surveys int
	var avgLatency float64
	err = h.DB.QueryRow(r.Context(), `SELECT
		(SELECT count(*) FROM page_views WHERE created_at >= $1 AND created_at < $2),
		(SELECT count(DISTINCT visitor) FROM page_views WHERE created_at >= $1 AND created_at < $2),
		(SELECT count(*) FROM chat_logs WHERE created_at >= $1 AND created_at < $2),
		(SELECT count(*) FILTER (WHERE refused) FROM chat_logs WHERE created_at >= $1 AND created_at < $2),
		(SELECT count(*) FILTER (WHERE helpful) FROM chat_logs WHERE created_at >= $1 AND created_at < $2),
		(SELECT count(*) FILTER (WHERE helpful = false) FROM chat_logs WHERE created_at >= $1 AND created_at < $2),
		(SELECT coalesce(avg(latency_ms),0) FROM chat_logs WHERE created_at >= $1 AND created_at < $2),
		(SELECT count(*) FROM survey_responses WHERE created_at >= $1 AND created_at < $2)`, from, to).
		Scan(&views, &visitors, &chats, &refused, &helpful, &unhelpful, &avgLatency, &surveys)
	if err != nil {
		return err
	}
	out["traffic"] = map[string]int{"views": views, "visitors": visitors}
	out["chat"] = map[string]any{"questions": chats, "refused": refused, "helpful": helpful, "unhelpful": unhelpful, "avg_latency_ms": int(avgLatency)}
	out["surveys"] = surveys

	daily, err := h.groupCount(r, `SELECT to_char(d, 'YYYY-MM-DD'), coalesce(c, 0)::int FROM generate_series($1::date, ($2::date - 1), '1 day') d
		LEFT JOIN (SELECT created_at::date AS dd, count(*) c FROM page_views WHERE created_at >= $1 AND created_at < $2 GROUP BY 1) v ON v.dd = d ORDER BY d`, from, to)
	if err != nil {
		return err
	}
	out["daily_views"] = daily
	if out["top_pages"], err = h.groupCount(r, `SELECT path, count(*)::int FROM page_views WHERE created_at >= $1 AND created_at < $2 GROUP BY 1 ORDER BY 2 DESC LIMIT 10`, from, to); err != nil {
		return err
	}
	if out["views_by_lang"], err = h.groupCount(r, `SELECT lang, count(*)::int FROM page_views WHERE created_at >= $1 AND created_at < $2 GROUP BY 1 ORDER BY 2 DESC`, from, to); err != nil {
		return err
	}
	if out["chat_by_lang"], err = h.groupCount(r, `SELECT lang, count(*)::int FROM chat_logs WHERE created_at >= $1 AND created_at < $2 GROUP BY 1 ORDER BY 2 DESC`, from, to); err != nil {
		return err
	}
	if out["chat_categories"], err = h.groupCount(r, `SELECT category, count(*)::int FROM chat_logs WHERE created_at >= $1 AND created_at < $2 GROUP BY 1 ORDER BY 2 DESC`, from, to); err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusOK, out)
}

type questionStat struct {
	ID      int64             `json:"id"`
	Section string            `json:"section"`
	Text    map[string]string `json:"text"`
	N       int               `json:"n"`
	Mean    float64           `json:"mean"`
	SD      float64           `json:"sd"`
	Level   string            `json:"level"`
	Dist    [5]int            `json:"distribution"`
}

func (h *Handler) surveyStats(r *http.Request) ([]questionStat, error) {
	from, to := dateRange(r)
	rows, err := h.DB.Query(r.Context(), `
		SELECT q.id, q.section, q.text, coalesce(array_agg((a.value)::int) FILTER (WHERE a.value IS NOT NULL), '{}')
		FROM survey_questions q
		LEFT JOIN survey_responses s ON s.created_at >= $1 AND s.created_at < $2
		LEFT JOIN LATERAL jsonb_each_text(s.answers) a ON a.key = q.id::text
		GROUP BY q.id ORDER BY q.section, q.sort, q.id`, from, to)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (questionStat, error) {
		var s questionStat
		var vals []int32
		if err := row.Scan(&s.ID, &s.Section, &s.Text, &vals); err != nil {
			return s, err
		}
		xs := make([]float64, len(vals))
		for i, v := range vals {
			xs[i] = float64(v)
			if v >= 1 && v <= 5 {
				s.Dist[v-1]++
			}
		}
		s.N = len(xs)
		s.Mean, s.SD = MeanSD(xs)
		s.Level = Level(s.Mean)
		return s, nil
	})
}

func (h *Handler) SurveyReport(w http.ResponseWriter, r *http.Request) error {
	stats, err := h.surveyStats(r)
	if err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusOK, stats)
}

// Export streams CSV (UTF-8 with BOM so Excel opens Thai/Chinese correctly) — scope 8.7 / 9.10.
func (h *Handler) Export(w http.ResponseWriter, r *http.Request) error {
	from, to := dateRange(r)
	kind := r.PathValue("kind")
	u := auth.FromContext(r.Context())
	// Raw survey comments and chat questions are personal research data: admins only.
	if (kind == "survey" || kind == "chat") && (u == nil || u.Role != "admin") {
		return httpx.Forbidden()
	}
	if u != nil {
		auth.Audit(r.Context(), h.DB, u.ID, "export", kind, 0, map[string]string{"from": from.Format("2006-01-02"), "to": to.Format("2006-01-02")})
	}
	var header []string
	var q string
	switch kind {
	case "survey":
		header = []string{"id", "created_at", "lang", "respondent", "answers", "comment"}
		q = `SELECT id::text, to_char(created_at,'YYYY-MM-DD HH24:MI'), lang, respondent::text, answers::text, comment FROM survey_responses WHERE created_at >= $1 AND created_at < $2 ORDER BY id`
	case "survey-summary":
		stats, err := h.surveyStats(r)
		if err != nil {
			return err
		}
		cw := startCSV(w, kind, []string{"question_id", "section", "question_th", "n", "mean", "sd", "level", "n1", "n2", "n3", "n4", "n5"})
		for _, s := range stats {
			cw.Write([]string{strconv.FormatInt(s.ID, 10), s.Section, httpx.CSVSafe(s.Text["th"]), strconv.Itoa(s.N),
				fmt.Sprintf("%.2f", s.Mean), fmt.Sprintf("%.2f", s.SD), s.Level,
				strconv.Itoa(s.Dist[0]), strconv.Itoa(s.Dist[1]), strconv.Itoa(s.Dist[2]), strconv.Itoa(s.Dist[3]), strconv.Itoa(s.Dist[4])})
		}
		cw.Flush()
		return cw.Error()
	case "chat":
		header = []string{"id", "created_at", "lang", "category", "question", "answer", "refused", "helpful", "latency_ms", "citations"}
		q = `SELECT id::text, to_char(created_at,'YYYY-MM-DD HH24:MI'), lang, category, question, answer, refused::text, coalesce(helpful::text,''), latency_ms::text, citations::text
			FROM chat_logs WHERE created_at >= $1 AND created_at < $2 ORDER BY id`
	case "views":
		header = []string{"day", "path", "lang", "views", "visitors"}
		q = `SELECT created_at::date::text, path, lang, count(*)::text, count(DISTINCT visitor)::text FROM page_views
			WHERE created_at >= $1 AND created_at < $2 GROUP BY 1,2,3 ORDER BY 1,4 DESC`
	default:
		return httpx.NotFound()
	}
	rows, err := h.DB.Query(r.Context(), q, from, to)
	if err != nil {
		return err
	}
	defer rows.Close()
	cw := startCSV(w, kind, header)
	for rows.Next() {
		vals := make([]string, len(header))
		ptrs := make([]any, len(header))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		for i, v := range vals {
			vals[i] = httpx.CSVSafe(v)
		}
		cw.Write(vals)
	}
	cw.Flush()
	return rows.Err()
}

func startCSV(w http.ResponseWriter, name string, header []string) *csv.Writer {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s-%s.csv"`, name, time.Now().Format("20060102")))
	w.Write([]byte("\xEF\xBB\xBF"))
	cw := csv.NewWriter(w)
	cw.Write(header)
	return cw
}

func (h *Handler) AuditLog(w http.ResponseWriter, r *http.Request) error {
	var out json.RawMessage
	err := h.DB.QueryRow(r.Context(), `SELECT coalesce(jsonb_agg(x ORDER BY x.id DESC), '[]') FROM (
		SELECT a.id, a.action, a.entity, a.entity_id, a.detail, a.created_at, u.email
		FROM audit_log a LEFT JOIN users u ON u.id = a.user_id ORDER BY a.id DESC LIMIT 300) x`).Scan(&out)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(out)
	return err
}

// Settings: small key/value store for project name, logo, contact, privacy text (scope 1.4).

var settingKeys = map[string]bool{"site": true, "privacy": true, "consent": true}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) error {
	var out json.RawMessage
	if err := h.DB.QueryRow(r.Context(), `SELECT coalesce(jsonb_object_agg(key, value), '{}') FROM site_settings`).Scan(&out); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(out)
	return err
}

func (h *Handler) PutSetting(w http.ResponseWriter, r *http.Request) error {
	key := r.PathValue("key")
	if !settingKeys[key] {
		return httpx.NotFound()
	}
	var v map[string]any
	if err := httpx.Decode(w, r, &v); err != nil {
		return err
	}
	b, _ := json.Marshal(v)
	if len(b) > 100_000 {
		return httpx.BadRequest("setting too large")
	}
	if _, err := h.DB.Exec(r.Context(), `INSERT INTO site_settings(key, value) VALUES ($1,$2) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, key, b); err != nil {
		return err
	}
	if u := auth.FromContext(r.Context()); u != nil {
		auth.Audit(r.Context(), h.DB, u.ID, "update", "setting:"+key, 0, nil)
	}
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}
