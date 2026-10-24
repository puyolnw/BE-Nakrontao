package httpx

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClientIPIgnoresSpoofedHeaderFromUntrustedPeer(t *testing.T) {
	SetTrustedProxies("172.16.0.0/12")
	defer SetTrustedProxies("")
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "203.0.113.9:5555"
	r.Header.Set("X-Forwarded-For", "1.2.3.4")
	if got := ClientIP(r); got != "203.0.113.9" {
		t.Fatalf("spoofed XFF accepted: %s", got)
	}
}

func TestClientIPUsesRightmostUntrustedHop(t *testing.T) {
	SetTrustedProxies("172.16.0.0/12")
	defer SetTrustedProxies("")
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "172.18.0.5:5555"
	r.Header.Set("X-Forwarded-For", "6.6.6.6, 198.51.100.7") // client prepended a fake hop
	if got := ClientIP(r); got != "198.51.100.7" {
		t.Fatalf("got %s", got)
	}
}

func TestBadProxyConfigRejected(t *testing.T) {
	if SetTrustedProxies("not-a-cidr") == nil {
		t.Fatal("expected error")
	}
	SetTrustedProxies("")
}

func TestLimiter(t *testing.T) {
	now := time.Unix(0, 0)
	l := NewLimiter(60, 2)
	l.now = func() time.Time { return now }
	if !l.Allow("a") || !l.Allow("a") || l.Allow("a") {
		t.Fatal("burst not enforced")
	}
	if !l.Allow("b") {
		t.Fatal("keys must be independent")
	}
	now = now.Add(time.Second)
	if !l.Allow("a") {
		t.Fatal("tokens must refill")
	}
}

func TestDailyCap(t *testing.T) {
	d := NewDailyCap(2)
	if !d.Take() || !d.Take() || d.Take() {
		t.Fatal("cap not enforced")
	}
}

func TestDecodeRejectsUnknownFieldsAndTrailingData(t *testing.T) {
	var v struct{ A int }
	for _, body := range []string{`{"A":1,"role":"admin"}`, `{"A":1}{"A":2}`, `not json`} {
		r := httptest.NewRequest("POST", "/", strings.NewReader(body))
		if Decode(httptest.NewRecorder(), r, &v) == nil {
			t.Errorf("accepted %s", body)
		}
	}
}

func TestDecodeRejectsOversizedBody(t *testing.T) {
	var v struct{ A string }
	big := `{"A":"` + strings.Repeat("x", maxBody+10) + `"}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(big))
	if Decode(httptest.NewRecorder(), r, &v) == nil {
		t.Fatal("oversized body accepted")
	}
}

func TestSecurityHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	SecurityHeaders(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, httptest.NewRequest("GET", "/api/x", nil))
	csp := w.Header().Get("Content-Security-Policy")
	if strings.Contains(csp, "script-src 'unsafe-inline'") || !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Fatalf("weak CSP: %s", csp)
	}
	for _, h := range []string{"X-Content-Type-Options", "Strict-Transport-Security", "X-Frame-Options"} {
		if w.Header().Get(h) == "" {
			t.Errorf("missing %s", h)
		}
	}
}

func TestShellCSPHashesInlineScript(t *testing.T) {
	f := t.TempDir() + "/index.html"
	os.WriteFile(f, []byte(`<html><script>boot()</script><script src="/a.js"></script></html>`), 0o600)
	csp := ShellCSP(f)
	// sha256("boot()") base64
	if !strings.Contains(csp, "'sha256-") || strings.Contains(csp, "'unsafe-inline' 'sha") || strings.Count(csp, "sha256-") != 1 {
		t.Fatalf("unexpected CSP: %s", csp)
	}
	if !strings.Contains(csp, "frame-ancestors 'none'") || !strings.Contains(csp, "object-src 'none'") {
		t.Fatal("base directives missing")
	}
}
