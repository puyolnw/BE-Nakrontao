package rag

import (
	"context"
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestChunkKeepsShortTextWhole(t *testing.T) {
	c := Chunk("เจ้าพ่อศรีนครเตา\nเป็นที่เคารพของชุมชน")
	if len(c) != 1 || !strings.Contains(c[0], "ชุมชน") {
		t.Fatalf("got %q", c)
	}
}

func TestChunkSplitsLongTextWithOverlapAndRuneLimit(t *testing.T) {
	para := strings.Repeat("ก", 400)
	text := strings.Join([]string{para, para, para, para}, "\n")
	c := Chunk(text)
	if len(c) < 2 {
		t.Fatalf("expected several chunks, got %d", len(c))
	}
	for _, s := range c {
		if utf8.RuneCountInString(s) > chunkRunes+chunkOverlap+1 {
			t.Fatalf("chunk too long: %d runes", utf8.RuneCountInString(s))
		}
		if !utf8.ValidString(s) {
			t.Fatal("chunk split inside a multi-byte rune")
		}
	}
}

func TestChunkEmpty(t *testing.T) {
	if Chunk("  \n ") != nil {
		t.Fatal("expected nil")
	}
}

func cos(a, b []float32) float64 {
	var s float64
	for i := range a {
		s += float64(a[i]) * float64(b[i])
	}
	return s
}

func TestHashEmbedderNormalisedAndSimilar(t *testing.T) {
	v, _ := HashEmbedder{}.Embed(context.Background(), []string{
		"ประวัติเจ้าพ่อศรีนครเตา", "เจ้าพ่อศรีนครเตามีประวัติอย่างไร", "ร้านอาหารใกล้ที่พัก",
	})
	var n float64
	for _, x := range v[0] {
		n += float64(x * x)
	}
	if math.Abs(n-1) > 1e-4 {
		t.Fatalf("not unit length: %f", n)
	}
	if cos(v[0], v[1]) <= cos(v[0], v[2]) {
		t.Fatal("related Thai sentences should be closer than unrelated ones")
	}
}

func TestRRFAndDedupe(t *testing.T) {
	a := []Hit{{ID: 1, DocCode: "A", VecSim: .9}, {ID: 2, DocCode: "B", VecSim: .5}}
	b := []Hit{{ID: 2, DocCode: "B"}, {ID: 3, DocCode: "C"}}
	f := RRF(a, b)
	if f[0].ID != 2 {
		t.Fatalf("chunk in both lists should rank first, got %+v", f)
	}
	if f[0].VecSim != .5 {
		t.Fatal("vector similarity must survive fusion for the refusal check")
	}
	d := dedupe([]Hit{{DocCode: "A", Lang: "th"}, {DocCode: "A", Lang: "th", KwSim: .9}, {DocCode: "A", Lang: "zh"}})
	if len(d) != 2 || d[0].KwSim != .9 {
		t.Fatalf("dedupe must merge and keep the strongest signal: %+v", d)
	}
	if got := DocRanking([]Hit{{DocCode: "A"}, {DocCode: "B"}, {DocCode: "A"}}); strings.Join(got, ",") != "A,B" {
		t.Fatal(got)
	}
}

func TestStrongEnough(t *testing.T) {
	if strongEnough(ModeHybrid, nil, .3) || strongEnough(ModeHybrid, []Hit{{VecSim: 0.01, KwSim: 0.2}}, .3) ||
		!strongEnough(ModeHybrid, []Hit{{VecSim: 0.5}}, .3) || !strongEnough(ModeHybrid, []Hit{{KwSim: 0.9}}, 0) {
		t.Fatal("threshold logic wrong")
	}
	// Offline embedder (minVec 0): semantic similarity alone must not pass in hybrid mode.
	if strongEnough(ModeHybrid, []Hit{{VecSim: 0.9, KwSim: 0.1}}, 0) {
		t.Fatal("hash embedder must require a keyword match")
	}
	if !strongEnough(ModeVector, []Hit{{VecSim: 0.4}}, 0) {
		t.Fatal("threshold logic wrong")
	}
}

func TestClassify(t *testing.T) {
	cases := map[string]string{
		"ต้องเตรียมของไหว้อะไรบ้าง": "ritual",
		"祭拜需要准备什么？":                  "ritual",
		"เจ้าพ่อมีประวัติอย่างไร":    "history",
		"How do I get there by car?": "travel",
		"สวัสดี":                     "other",
	}
	for q, want := range cases {
		if got := Classify(q); got != want {
			t.Errorf("%s: got %s want %s", q, got, want)
		}
	}
}

func TestExtractiveReturnsFirstPassage(t *testing.T) {
	out, _ := Extractive{}.Complete(context.Background(), "", buildPrompt("th", "q", []Hit{{DocCode: "K1", Content: "หนึ่ง"}, {DocCode: "K2", Content: "สอง"}}))
	if out != "หนึ่ง" {
		t.Fatalf("got %q", out)
	}
}
