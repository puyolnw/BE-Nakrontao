package httpx

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// Limiter is an in-memory token bucket per key. One process serves the whole site, so no Redis is needed.
type Limiter struct {
	mu      sync.Mutex
	rate    float64 // tokens per second
	burst   float64
	buckets map[string]*bucket
	now     func() time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

func NewLimiter(perMinute, burst int) *Limiter {
	return &Limiter{rate: float64(perMinute) / 60, burst: float64(burst), buckets: map[string]*bucket{}, now: time.Now}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	b, ok := l.buckets[key]
	if !ok {
		if len(l.buckets) > 50_000 {
			l.sweep(now)
		}
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	b.tokens = min(l.burst, b.tokens+now.Sub(b.last).Seconds()*l.rate)
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (l *Limiter) sweep(now time.Time) {
	for k, b := range l.buckets {
		if now.Sub(b.last) > 10*time.Minute {
			delete(l.buckets, k)
		}
	}
}

// DailyCap counts events per UTC day and refuses after max.
type DailyCap struct {
	mu    sync.Mutex
	max   int
	day   string
	count int
}

func NewDailyCap(max int) *DailyCap { return &DailyCap{max: max} }

// KeyedDailyCap allows max events per key per UTC day, e.g. per visitor. Counts reset at midnight UTC.
type KeyedDailyCap struct {
	mu     sync.Mutex
	max    int
	day    string
	counts map[string]int
}

func NewKeyedDailyCap(max int) *KeyedDailyCap {
	return &KeyedDailyCap{max: max, counts: map[string]int{}}
}

func (d *KeyedDailyCap) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	today := time.Now().UTC().Format("2006-01-02")
	if today != d.day {
		d.day, d.counts = today, map[string]int{}
	}
	if d.counts[key] >= d.max {
		return false
	}
	d.counts[key]++
	return true
}

func (d *DailyCap) Take() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	today := time.Now().UTC().Format("2006-01-02")
	if today != d.day {
		d.day, d.count = today, 0
	}
	if d.count >= d.max {
		return false
	}
	d.count++
	return true
}

// APILimit caps every /api request per client address (scraping, brute force, cheap DoS).
// Reads and writes have separate budgets; anything over gets 429 before touching the database.
func APILimit(reads, writes *Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") && r.URL.Path != "/api/health" {
			l := reads
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				l = writes
			}
			if !l.Allow(ClientIP(r)) {
				w.Header().Set("Retry-After", "30")
				WriteError(w, r, TooMany())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
