package rate_limiter

// interface para implementação do strategy
type Strategy interface {
	Allow(ip string, token string) bool
}
