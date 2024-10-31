package rate_limiter

type RateLimiter interface {
	Allow(ip string, token string) bool
}

type rateLimiter struct {
	strategy Strategy
}

func NewRateLimiter(strategy Strategy) RateLimiter {
	return &rateLimiter{strategy: strategy}
}

func (r *rateLimiter) Allow(ip string, token string) bool {
	return r.strategy.Allow(ip, token)
}
