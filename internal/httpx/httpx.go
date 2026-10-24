// Package httpx holds JSON helpers, error responses and shared middleware.
package httpx

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const maxBody = 1 << 20

type Error struct {
	Status  int
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

func BadRequest(msg string) *Error { return &Error{http.StatusBadRequest, "bad_request", msg} }
func NotFound() *Error             { return &Error{http.StatusNotFound, "not_found", "resource not found"} }
func Unauthorized() *Error         { return &Error{http.StatusUnauthorized, "unauthorized", "login required"} }
func Forbidden() *Error            { return &Error{http.StatusForbidden, "forbidden", "insufficient role"} }
func TooMany() *Error              { return &Error{http.StatusTooManyRequests, "rate_limited", "too many requests"} }

// Handler is an http handler that returns an error; Wrap turns it into an http.HandlerFunc.
type Handler func(w http.ResponseWriter, r *http.Request) error

func Wrap(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			WriteError(w, r, err)
		}
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var he *Error
	if !errors.As(err, &he) {
		slog.Error("request failed", "path", r.URL.Path, "err", err)
		he = &Error{http.StatusInternalServerError, "internal", "internal server error"}
	}
	JSON(w, he.Status, map[string]string{"error": he.Code, "message": he.Message})
}

func JSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Decode reads a size-limited JSON body and rejects unknown fields.
func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return BadRequest("invalid JSON: " + err.Error())
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		return BadRequest("body must contain a single JSON object")
	}
	return nil
}

func PathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, BadRequest("invalid id")
	}
	return id, nil
}

// Lang returns a supported language code from ?lang=, defaulting to Thai.
func Lang(r *http.Request) string {
	switch l := r.URL.Query().Get("lang"); l {
	case "th", "zh", "en":
		return l
	}
	return "th"
}

// TrustedProxies lists networks whose X-Forwarded-For is believed (the reverse proxy). Set once at startup.
var TrustedProxies []*net.IPNet

func SetTrustedProxies(cidrs string) error {
	TrustedProxies = nil
	for _, c := range strings.Split(cidrs, ",") {
		if c = strings.TrimSpace(c); c == "" {
			continue
		}
		_, n, err := net.ParseCIDR(c)
		if err != nil {
			return fmt.Errorf("TRUSTED_PROXY %q: %w", c, err)
		}
		TrustedProxies = append(TrustedProxies, n)
	}
	return nil
}

func trusted(ip net.IP) bool {
	for _, n := range TrustedProxies {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// ClientIP returns the peer address. X-Forwarded-For is only honoured when the peer is a trusted
// proxy, and then the right-most untrusted hop is used, so clients cannot spoof their IP.
func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(host)
	if ip == nil || !trusted(ip) {
		return host
	}
	hops := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	for i := len(hops) - 1; i >= 0; i-- {
		hop := net.ParseIP(strings.TrimSpace(hops[i]))
		if hop == nil {
			break
		}
		if !trusted(hop) {
			return hop.String()
		}
	}
	return host
}

// HashIP pseudonymises an address so analytics never store raw IPs.
func HashIP(salt, ip string) string {
	sum := sha256.Sum256([]byte(salt + "|" + ip + "|" + time.Now().UTC().Format("2006-01-02")))
	return hex.EncodeToString(sum[:12])
}

// CSVSafe neutralises spreadsheet formula injection in exported cells.
func CSVSafe(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}
	return s
}

const baseCSP = "default-src 'self'; img-src 'self' data: https:; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; " +
	"frame-src https://www.youtube-nocookie.com; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'; " +
	"upgrade-insecure-requests"

var inlineScript = regexp.MustCompile(`(?s)<script>(.*?)</script>`)

// ShellCSP returns the page CSP with a hash for every inline <script> in the HTML file, so the
// SvelteKit bootstrap runs while injected inline scripts are still blocked.
func ShellCSP(htmlFile string) string {
	src := "script-src 'self'"
	if b, err := os.ReadFile(htmlFile); err == nil {
		for _, m := range inlineScript.FindAllSubmatch(b, -1) {
			sum := sha256.Sum256(m[1])
			src += " 'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
		}
	}
	return baseCSP + "; " + src
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.Set("Permissions-Policy", "geolocation=(self), camera=(), microphone=(), payment=(), usb=()")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		h.Set("X-Permitted-Cross-Domain-Policies", "none")
		h.Set("Content-Security-Policy", baseCSP+"; script-src 'self'")
		next.ServeHTTP(w, r)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(c int) { s.status = c; s.ResponseWriter.WriteHeader(c) }

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{w, 200}
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic", "path", r.URL.Path, "err", fmt.Sprint(rec))
				WriteError(sw, r, errors.New("panic"))
			}
			if strings.HasPrefix(r.URL.Path, "/api/") {
				slog.Info("http", "method", r.Method, "path", r.URL.Path, "status", sw.status, "ms", time.Since(start).Milliseconds())
			}
		}()
		next.ServeHTTP(sw, r)
	})
}
