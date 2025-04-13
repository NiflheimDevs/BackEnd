package midratelimit

import (
	"net/http"

	"github.com/niflheimdevs/backend/internal/domain/exceptions"
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
			err := exceptions.Exception{
				Tag:    exceptions.NOT_FOUND,
				Errors: []exceptions.SpecificError{},
			}
			panic(err)
		}
		next.ServeHTTP(w, r)
	})
}
