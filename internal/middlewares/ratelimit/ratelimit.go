package midratelimit

import (
	"net/http"

	"golang.org/x/time/rate"
)

type RateLimit struct {
	limit rate.Limit
	burst int
}

func NewRateLimit() *RateLimit {
	return &RateLimit{
		limit: 5,
		burst: 10,
	}
}

func (rl *RateLimit) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := rate.NewLimiter(rl.limit, rl.burst)
		if !limiter.Allow() {
			panic("Rate Limit Exceed")
		}
		next.ServeHTTP(w, r)
	})
}
