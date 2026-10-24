package rag

import (
	"context"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Mode selects the retrieval strategy. The four modes are the research conditions (FF proposal §4.1).
type Mode string

const (
	ModeNone         Mode = "llm_only"     // A: baseline, no retrieval
	ModeVector       Mode = "vector"       // B: standard RAG, same-language vector search
	ModeMultilingual Mode = "multilingual" // C: vector search across all languages
	ModeHybrid       Mode = "hybrid"       // D: vector + keyword, RRF fusion, lexical rerank (production)
)

func ValidMode(m string) bool {
	switch Mode(m) {
	case ModeNone, ModeVector, ModeMultilingual, ModeHybrid:
		return true
	}
	return false
}

type Hit struct {
	ID      int64   `json:"-"`
	DocCode string  `json:"doc_code"`
	Lang    string  `json:"lang"`
	Kind    string  `json:"kind"`
	Title   string  `json:"title"`
	Content string  `json:"-"`
	Source  string  `json:"source"`
	Score   float64 `json:"score"`
	VecSim  float64 `json:"-"` // best cosine similarity from vector search
	KwSim   float64 `json:"-"` // best trigram word similarity from keyword search
}

type Retriever struct {
	DB       *pgxpool.Pool
	Embedder Embedder
}

const (
	candidatePool = 30
	// keywordFloor admits weak keyword matches into fusion; the evidence gate decides separately.
	// The knowledge base is a few hundred chunks, so a sequential scan here costs well under 10 ms.
	keywordFloor = 0.25
)

func (rt *Retriever) Retrieve(ctx context.Context, mode Mode, lang, query string, k int) ([]Hit, error) {
	if mode == ModeNone {
		return nil, nil
	}
	vecs, err := rt.Embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	vec := vectorLiteral(vecs[0])
	langFilter := ""
	if mode == ModeVector {
		langFilter = lang
	}
	vecHits, err := rt.vectorSearch(ctx, vec, langFilter)
	if err != nil {
		return nil, err
	}
	if mode != ModeHybrid {
		return topK(dedupe(vecHits), k), nil
	}
	kwHits, err := rt.keywordSearch(ctx, query)
	if err != nil {
		return nil, err
	}
	fused := RRF(vecHits, kwHits)
	Rerank(query, fused)
	return topK(dedupe(fused), k), nil
}

func (rt *Retriever) vectorSearch(ctx context.Context, vec, lang string) ([]Hit, error) {
	rows, err := rt.DB.Query(ctx, `
		SELECT id, doc_code, lang, kind, title, content, source, 1 - (embedding <=> $1::vector) AS sim
		FROM kb_chunks WHERE ($2 = '' OR lang = $2)
		ORDER BY embedding <=> $1::vector LIMIT $3`, vec, lang, candidatePool)
	if err != nil {
		return nil, err
	}
	return collect(rows, true)
}

// keywordSearch uses trigram word similarity, which works on Thai/Chinese because it needs no tokenizer.
func (rt *Retriever) keywordSearch(ctx context.Context, q string) ([]Hit, error) {
	rows, err := rt.DB.Query(ctx, `
		SELECT id, doc_code, lang, kind, title, content, source, word_similarity($1, title || ' ' || content) AS sim
		FROM kb_chunks WHERE word_similarity($1, title || ' ' || content) >= $2
		ORDER BY sim DESC LIMIT $3`, q, keywordFloor, candidatePool)
	if err != nil {
		return nil, err
	}
	return collect(rows, false)
}

func collect(rows pgx.Rows, vector bool) ([]Hit, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Hit, error) {
		var h Hit
		err := row.Scan(&h.ID, &h.DocCode, &h.Lang, &h.Kind, &h.Title, &h.Content, &h.Source, &h.Score)
		if vector {
			h.VecSim = h.Score
		} else {
			h.KwSim = h.Score
		}
		return h, err
	})
}

// RRF fuses ranked lists with Reciprocal Rank Fusion (k=60), keyed by chunk id.
func RRF(lists ...[]Hit) []Hit {
	const k = 60.0
	byID := map[int64]*Hit{}
	var order []int64
	for _, list := range lists {
		for rank, h := range list {
			cur, ok := byID[h.ID]
			if !ok {
				c := h
				c.Score = 0
				byID[h.ID] = &c
				cur = &c
				order = append(order, h.ID)
			}
			cur.VecSim = max(cur.VecSim, h.VecSim)
			cur.KwSim = max(cur.KwSim, h.KwSim)
			cur.Score += 1 / (k + float64(rank+1))
		}
	}
	out := make([]Hit, 0, len(order))
	for _, id := range order {
		out = append(out, *byID[id])
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

// Rerank adds a lexical overlap bonus (character bigrams of query found in title+content) and re-sorts.
func Rerank(query string, hits []Hit) {
	qg := bigrams(query)
	if len(qg) == 0 {
		return
	}
	for i := range hits {
		dg := bigrams(hits[i].Title + " " + hits[i].Content)
		match := 0
		for g := range qg {
			if dg[g] {
				match++
			}
		}
		// Small bonuses break RRF ties in favour of chunks that are actually close to the query.
		hits[i].Score += 0.02*float64(match)/float64(len(qg)) + 0.02*hits[i].KwSim + 0.02*hits[i].VecSim
	}
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
}

func bigrams(s string) map[string]bool {
	rs := normalizeText(s)
	m := map[string]bool{}
	for i := 0; i+2 <= len(rs); i++ {
		if rs[i] != ' ' && rs[i+1] != ' ' {
			m[string(rs[i:i+2])] = true
		}
	}
	return m
}

// dedupe keeps the best-ranked chunk per document and language so citations are not repeated,
// carrying over the strongest similarity signals of the dropped chunks for the evidence gate.
func dedupe(hits []Hit) []Hit {
	idx := map[string]int{}
	var out []Hit
	for _, h := range hits {
		key := h.DocCode + "|" + h.Lang
		if i, ok := idx[key]; ok {
			out[i].VecSim = max(out[i].VecSim, h.VecSim)
			out[i].KwSim = max(out[i].KwSim, h.KwSim)
			continue
		}
		idx[key] = len(out)
		out = append(out, h)
	}
	return out
}

func topK(h []Hit, k int) []Hit {
	if len(h) > k {
		return h[:k]
	}
	return h
}

// DocRanking returns distinct doc codes in rank order (for Recall@K / MRR).
func DocRanking(hits []Hit) []string {
	seen := map[string]bool{}
	var out []string
	for _, h := range hits {
		if !seen[h.DocCode] {
			seen[h.DocCode] = true
			out = append(out, h.DocCode)
		}
	}
	return out
}
