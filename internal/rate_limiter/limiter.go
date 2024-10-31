package rate_limiter

// Inteface para implmentação do RateLimiter
type RateLimiter interface {
	Allow(ip string, token string) bool
}

type rateLimiter struct {
	strategy Strategy
}

// Inicialização do objeto RateLimiter
func NewRateLimiter(strategy Strategy) RateLimiter {
	return &rateLimiter{strategy: strategy}
}

// Função de validação das conexões
func (r *rateLimiter) Allow(ip string, token string) bool {
	return r.strategy.Allow(ip, token)
}
