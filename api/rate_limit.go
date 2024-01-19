package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type limiter struct {
	count int
	iat   int64
}

type rateLimiter struct {
	lock     sync.Mutex
	limit    int
	lifetime time.Duration
	limits   map[string]*limiter
}

func NewRateLimiter(limit int, lifetime time.Duration) *rateLimiter {
	return &rateLimiter{
		limits:   make(map[string]*limiter),
		limit:    limit,
		lifetime: lifetime,
	}
}

func (r *rateLimiter) checkLimit(key string) (int, int64) {
	r.lock.Lock()
	defer r.lock.Unlock()

	lim, ok := r.limits[key]
	if !ok {
		// create new limiter
		lim = &limiter{
			count: 0,
			iat:   time.Now().Unix(),
		}
		r.limits[key] = lim
		return r.limit, lim.iat + int64(r.lifetime.Seconds())
	}

	reset := lim.iat + int64(r.lifetime.Seconds())

	// if the limiter is expired, reset it
	if time.Unix(lim.iat, 0).Add(r.lifetime).Before(time.Now()) {
		lim.count = 0
		lim.iat = time.Now().Unix()
		return r.limit, reset
	}

	// if rate limit is reached, return false
	if lim.count >= r.limit {
		return 0, reset
	}

	// incr the counter
	lim.count++
	return r.limit - lim.count, reset
}

type RateLimitStrategy string

const (
	RateLimitStrategyIP   RateLimitStrategy = "ip"
	RateLimitStrategyKey  RateLimitStrategy = "key"
	RateLimitStrategyNone RateLimitStrategy = "none"
)

func ToRateLimitStrategy(s string) RateLimitStrategy {
	switch s {
	case "ip":
		return RateLimitStrategyIP
	case "key":
		return RateLimitStrategyKey
	case "none", "", "false":
		return RateLimitStrategyNone
	default:
		slog.Warn("invalid rate limit strategy, using no strategy", "rate_limit_strategy", s)
		return RateLimitStrategyNone
	}
}

func (s *Server) rateLimitMiddleware(limit int, stratrgy string) func(next http.Handler) http.Handler {
	strat := ToRateLimitStrategy(stratrgy)

	if strat == RateLimitStrategyNone || limit <= 0 || s.rateLimit == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			var k string

			switch strat {
			case RateLimitStrategyIP:
				k = r.RemoteAddr

			case RateLimitStrategyKey:
				authHeader := r.Header.Get("Authentication")
				authHeader = strings.TrimPrefix(authHeader, "Bearer ")
				if authHeader == "" {
					writeError(w, http.StatusUnauthorized, errMissingAuth)
					return
				}
				k = authHeader
			}

			remaining, resetTime := s.rateLimit.checkLimit(k)
			if remaining <= 0 {
				writeJSON(w, http.StatusTooManyRequests, &JRPCResponse{
					JSONRPC: "2.0",
					Error: &JRPCError{
						Code:    -32900,
						Message: "Too Many Requests",
					},
				})
				return
			}
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}
}
