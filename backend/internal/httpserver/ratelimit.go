package httpserver

import (
	"net"
	"net/http"
	"sync"
	"time"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/httputil"
)

// rateLimiter is a fixed-window counter per client IP. State is in-process on
// purpose: one binary and no Redis, so a restart forgives every caller.
type rateLimiter struct {
	limit  int
	window time.Duration

	mu        sync.Mutex
	hits      map[string]*rateWindow
	lastSweep time.Time
}

type rateWindow struct {
	count   int
	resetAt time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:     limit,
		window:    window,
		hits:      map[string]*rateWindow{},
		lastSweep: time.Now(),
	}
}

func (l *rateLimiter) allow(key string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.lastSweep) > l.window {
		for k, w := range l.hits {
			if now.After(w.resetAt) {
				delete(l.hits, k)
			}
		}
		l.lastSweep = now
	}

	w, ok := l.hits[key]
	if !ok || now.After(w.resetAt) {
		l.hits[key] = &rateWindow{count: 1, resetAt: now.Add(l.window)}
		return true
	}
	if w.count >= l.limit {
		return false
	}
	w.count++
	return true
}

func (l *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.limit > 0 && !l.allow(clientIP(r)) {
			httputil.Error(w, apierr.RateLimited("too many attempts, try again later"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
