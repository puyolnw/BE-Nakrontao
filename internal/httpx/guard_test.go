package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeRequiresJSON(t *testing.T) {
	cases := map[string]int{
		"application/json":                  0,
		"application/json; charset=utf-8":   0,
		"text/plain":                        http.StatusUnsupportedMediaType,
		"application/x-www-form-urlencoded": http.StatusUnsupportedMediaType,
		"":                                  http.StatusUnsupportedMediaType,
	}
	for ct, want := range cases {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"a":1}`))
		if ct != "" {
			r.Header.Set("Content-Type", ct)
		}
		var v struct {
			A int `json:"a"`
		}
		err := Decode(httptest.NewRecorder(), r, &v)
		got := 0
		if he, ok := err.(*Error); ok && he.Status == http.StatusUnsupportedMediaType {
			got = he.Status
		}
		if got != want {
			t.Errorf("Content-Type %q: got status %d, want %d (err=%v)", ct, got, want, err)
		}
	}
}

func TestKeyedDailyCapPerKey(t *testing.T) {
	d := NewKeyedDailyCap(2)
	if !d.Allow("a") || !d.Allow("a") || d.Allow("a") {
		t.Fatal("key a must stop after 2")
	}
	if !d.Allow("b") {
		t.Fatal("key b must have its own allowance")
	}
}
