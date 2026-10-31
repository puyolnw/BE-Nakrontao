package rag

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDetectLang(t *testing.T) {
	cases := []struct{ q, page, want string }{
		{"Who is Chao Pho?", "th", "en"},
		{"师那空涛神是谁？", "en", "zh"},
		{"เจ้าพ่อคือใคร", "zh", "th"},
		{"K-001 คืออะไร", "en", "th"},
		{"Si Nakhon Tao 是谁", "th", "zh"},
		{"12345 ???", "zh", "zh"},
	}
	for _, c := range cases {
		if got := detectLang(c.q, c.page); got != c.want {
			t.Errorf("detectLang(%q, %q) = %q, want %q", c.q, c.page, got, c.want)
		}
	}
}

func TestScreen(t *testing.T) {
	cases := []struct {
		q    string
		want bool
	}{
		{"วิธีทำ ระเบิด ที่บ้าน", true},
		{"how to BUILD a bomb", true},
		{"如何洗钱", true},
		{"บั้งไฟระเบิดไหม", false},
		{"ประวัติเจ้าพ่อศรีนครเตา", false},
	}
	for _, c := range cases {
		if _, got := screen("th", c.q); got != c.want {
			t.Errorf("screen(%q) = %v, want %v", c.q, got, c.want)
		}
	}
}

func TestScreenCatchesCrisisWithPunctuationAndSpacing(t *testing.T) {
	for _, q := range []string{
		"ไม่ไหวแล้ว อยากตาย",
		"อยาก-ตาย",
		"อยาก ๆ จบ ชีวิต",
		"ไม่อยากอยู่แล้ว",
		"I want to end it all",
		"i don't want to live",
		"活不下去了",
		"想 死",
		"I WANT TO KILL MY-SELF",
		"อยากตาย​", // zero-width space
	} {
		reply, ok := screen("th", q)
		if !ok || !strings.Contains(reply, "1323") {
			t.Errorf("screen(%q) = %q, %v; want the hotline reply", q, reply, ok)
		}
	}
}

func TestCrisisBeatsForbiddenTerm(t *testing.T) {
	reply, ok := screen("th", "อยากตาย และอยากทำระเบิด")
	if !ok || !strings.Contains(reply, "1323") {
		t.Fatalf("crisis must win over forbidden terms, got %q", reply)
	}
}

func TestShapeCrisisToken(t *testing.T) {
	a := shape("en", "question", Answer{Answer: crisisToken, Citations: []Hit{{DocCode: "K-001"}}})
	if !a.Refused || len(a.Citations) != 0 || a.Answer != crisisReplies["en"] {
		t.Fatalf("crisis token not turned into the help reply: %+v", a)
	}
}

func TestShapePlayDumbToken(t *testing.T) {
	a := shape("th", "คำถาม", Answer{Answer: playDumbToken, Citations: []Hit{{DocCode: "K-001"}}})
	if !a.Refused || len(a.Citations) != 0 || strings.Contains(a.Answer, playDumbToken) {
		t.Fatalf("token leaked or not refused: %+v", a)
	}
}

func TestShapeReplacesUngroundedAnswers(t *testing.T) {
	cases := map[string]Answer{
		"refused by evidence gate":    {Answer: refusals["th"], Refused: true, Citations: []Hit{{DocCode: "K-001"}}},
		"answer without any citation": {Answer: "ทายผลบอลได้ว่าทีมไหนชนะ", Citations: []Hit{}},
		"answer citing nothing kept":  {Answer: "มีโทรศัพท์ไหม", Citations: nil},
	}
	for name, in := range cases {
		a := shape("th", "ฟุตบอลใครชนะ", in)
		if a.Answer != redirects["th"] || !a.Refused || len(a.Citations) != 0 {
			t.Errorf("%s: got %+v", name, a)
		}
	}
}

func TestShapeKeepsCitedAnswers(t *testing.T) {
	in := Answer{Answer: "ศาลอยู่บ้านเมืองเตา [K-001]", Citations: []Hit{{DocCode: "K-001"}}}
	if a := shape("th", "ศาลอยู่ไหน", in); a.Answer != in.Answer || a.Refused || len(a.Citations) != 1 {
		t.Fatalf("cited answer changed: %+v", a)
	}
}

func TestOpenAICompat(t *testing.T) {
	var got struct {
		Model    string `json:"model"`
		Messages []struct {
			Role string `json:"role"`
		} `json:"messages"`
	}
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"<think>hmm</think> สวัสดี "}}]}`))
	}))
	defer srv.Close()
	o := &OpenAICompat{URL: srv.URL, Key: "k", Model: "m"}
	text, err := o.Complete(context.Background(), "sys", "usr")
	if err != nil || text != "สวัสดี" {
		t.Fatalf("text = %q, err = %v", text, err)
	}
	if auth != "Bearer k" || got.Model != "m" || len(got.Messages) != 2 || got.Messages[0].Role != "system" {
		t.Fatalf("bad request: auth=%q body=%+v", auth, got)
	}
}

func TestOpenAICompatQuotaIsBusy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	_, err := (&OpenAICompat{URL: srv.URL}).Complete(context.Background(), "s", "u")
	if !errors.Is(err, error(ErrBusy)) {
		t.Fatalf("err = %v, want ErrBusy", err)
	}
}
