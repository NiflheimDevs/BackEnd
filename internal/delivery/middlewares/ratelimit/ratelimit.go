package midratelimit

import (
	"net/http"
	"sync"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiters map[string]*rate.Limiter
	limit    rate.Limit
	burst    int
	mu       sync.Mutex
}

func NewRateLimit(Constants *bootstrap.Constants) *RateLimit {
	return &RateLimit{
		limiters: make(map[string]*rate.Limiter),
		limit:    rate.Limit(Constants.RateLimiter.Limit),
		burst:    Constants.RateLimiter.Burst,
	}
}

func (rl *RateLimit) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.limit, rl.burst)
		rl.limiters[ip] = limiter
	}
	return limiter
}

func (rl *RateLimit) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			err := exceptions.Exception{
				Tag: exceptions.LIMIT_EXCEED,
			}
			panic(err)
		}
		next.ServeHTTP(w, r)
	})
}
