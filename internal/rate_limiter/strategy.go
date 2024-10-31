package rate_limiter

type Strategy interface {
	Allow(ip string, token string) bool
}
