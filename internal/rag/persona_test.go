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

type fakeLLM struct {
	reply string
	err   error
}

func (f fakeLLM) Name() string { return "fake" }

func (f fakeLLM) Complete(context.Context, string, string) (string, error) { return f.reply, f.err }

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

func TestScreenGivesHelpNotJokesInCrisis(t *testing.T) {
	reply, ok := screen("th", "ไม่ไหวแล้ว อยากตาย")
	if !ok || !strings.Contains(reply, "1323") {
		t.Fatalf("crisis reply = %q, %v", reply, ok)
	}
}

func TestShapeFlaggedByModel(t *testing.T) {
	s := &Service{LLM: fakeLLM{}}
	a := s.shape(context.Background(), "th", "คำถาม", Answer{Answer: playDumbToken, Citations: []Hit{{DocCode: "K-001"}}})
	if !a.Refused || len(a.Citations) != 0 || strings.Contains(a.Answer, playDumbToken) {
		t.Fatalf("token leaked or not refused: %+v", a)
	}
}

func TestShapeBridgesOffTopic(t *testing.T) {
	ctx := context.Background()
	refused := Answer{Answer: refusals["th"], Refused: true, Citations: []Hit{{DocCode: "K-001"}}}
	cases := []struct {
		name string
		llm  LLM
		want string
	}{
		{"model links back", fakeLLM{reply: "ลองชิมข้าวหอมมะลิ"}, "ลองชิมข้าวหอมมะลิ"},
		{"no llm", Extractive{}, redirects["th"]},
		{"provider down", fakeLLM{err: errors.New("down")}, redirects["th"]},
		{"no answer token", fakeLLM{reply: "NO_ANSWER"}, redirects["th"]},
	}
	for _, c := range cases {
		a := (&Service{LLM: c.llm}).shape(ctx, "th", "ฟุตบอลใครชนะ", refused)
		if a.Answer != c.want || !a.Refused || len(a.Citations) != 0 {
			t.Errorf("%s: got %+v", c.name, a)
		}
	}
}

func TestShapeKeepsRealAnswers(t *testing.T) {
	s := &Service{LLM: fakeLLM{reply: "unused"}}
	in := Answer{Answer: "ศาลอยู่บ้านเมืองเตา [K-001]", Citations: []Hit{{DocCode: "K-001"}}}
	if a := s.shape(context.Background(), "th", "ศาลอยู่ไหน", in); a.Answer != in.Answer || len(a.Citations) != 1 {
		t.Fatalf("answer changed: %+v", a)
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
