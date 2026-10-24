// Package rag implements the multilingual retrieval-augmented chatbot (scope 5.x) and the
// retrieval modes compared in the research experiments (R2–R5).
package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const Dim = 256

type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	Name() string
}

// HashEmbedder maps character n-grams into a fixed vector. Thai and Chinese have no spaces, so
// character n-grams work where word tokens would not. It needs no network and is deterministic,
// which makes it the offline default and a reproducible baseline for experiments.
type HashEmbedder struct{}

func (HashEmbedder) Name() string { return "hash-ngram-256" }

func (HashEmbedder) Embed(_ context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i, t := range texts {
		out[i] = hashVector(t)
	}
	return out, nil
}

func normalizeText(s string) []rune {
	var rs []rune
	space := false
	for _, r := range strings.ToLower(s) {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if !space && len(rs) > 0 {
				rs = append(rs, ' ')
			}
			space = true
			continue
		}
		space = false
		rs = append(rs, r)
	}
	return rs
}

func hashVector(s string) []float32 {
	v := make([]float32, Dim)
	add := func(feature string, w float32) {
		h := fnv.New32a()
		h.Write([]byte(feature))
		x := h.Sum32()
		sign := float32(1)
		if x&1 == 1 {
			sign = -1
		}
		v[(x>>1)%Dim] += sign * w
	}
	for _, run := range scriptRuns(s) {
		rs := []rune(run.text)
		switch run.script {
		case scriptLatin:
			// Whole words carry the meaning; in-word trigrams tolerate typos and inflection.
			if latinStop[run.text] {
				continue
			}
			add("w:"+run.text, 3)
			padded := []rune("<" + run.text + ">")
			for i := 0; i+3 <= len(padded); i++ {
				add("t:"+string(padded[i:i+3]), 1)
			}
		case scriptHan:
			// Chinese: single characters are often words; bigrams capture compounds.
			for i := range rs {
				add("h1:"+string(rs[i]), 1)
				if i+2 <= len(rs) {
					add("h2:"+string(rs[i:i+2]), 2)
				}
			}
		default:
			// Thai has no spaces and single letters mean nothing: use bigrams and trigrams.
			for i := 0; i+2 <= len(rs); i++ {
				add("b:"+string(rs[i:i+2]), 1)
				if i+3 <= len(rs) {
					add("c:"+string(rs[i:i+3]), 2)
				}
			}
		}
	}
	return l2(v)
}

const (
	scriptLatin = iota
	scriptHan
	scriptOther
)

type textRun struct {
	script int
	text   string
}

// scriptRuns lowercases s and splits it into runs of one script, dropping spaces and punctuation.
func scriptRuns(s string) []textRun {
	var runs []textRun
	var cur []rune
	curScript := -1
	flush := func() {
		if len(cur) > 0 {
			runs = append(runs, textRun{curScript, string(cur)})
		}
		cur = nil
	}
	for _, r := range strings.ToLower(s) {
		if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			flush()
			continue
		}
		sc := scriptOther
		switch {
		case unicode.Is(unicode.Han, r):
			sc = scriptHan
		case r < 0x250: // Latin letters and ASCII digits
			sc = scriptLatin
		}
		if sc != curScript {
			flush()
			curScript = sc
		}
		cur = append(cur, r)
	}
	flush()
	return runs
}

var latinStop = map[string]bool{}

func init() {
	for _, w := range strings.Fields(`a an the and or of to in on at for with by from is are was were be been am do does did i me my you your we our it its this that these those what which who whom how why when where can could should would will shall may might must there here please tell about some any`) {
		latinStop[w] = true
	}
}

func l2(v []float32) []float32 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	if sum == 0 {
		return v
	}
	n := float32(1 / math.Sqrt(sum))
	for i := range v {
		v[i] *= n
	}
	return v
}

// OpenAIEmbedder calls any OpenAI-compatible /v1/embeddings endpoint with dimensions=256.
type OpenAIEmbedder struct {
	URL, Key, Model string
	Client          *http.Client
}

func (e *OpenAIEmbedder) Name() string { return e.Model }

func (e *OpenAIEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	body, _ := json.Marshal(map[string]any{"model": e.Model, "input": texts, "dimensions": Dim})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.Key)
	c := e.Client
	if c == nil {
		c = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding api: status %d", resp.StatusCode)
	}
	var out struct {
		Data []struct {
			Index     int       `json:"index"`
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	vecs := make([][]float32, len(texts))
	for _, d := range out.Data {
		if d.Index < len(vecs) && len(d.Embedding) == Dim {
			vecs[d.Index] = l2(d.Embedding)
		}
	}
	for i, v := range vecs {
		if v == nil {
			return nil, fmt.Errorf("embedding api: missing vector %d", i)
		}
	}
	return vecs, nil
}

// vectorLiteral formats a vector for pgvector's text input.
func vectorLiteral(v []float32) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, x := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%.6f", x)
	}
	b.WriteByte(']')
	return b.String()
}
