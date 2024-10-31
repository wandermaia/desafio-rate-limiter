package rate_limiter

import (
	"time"

	"github.com/wandermaia/desafio-rate-limiter/internal/repository"
)

// Objeto do strategy redis
type redisStrategy struct {
	repo               *repository.RedisRepository
	maxRequests        int
	blockDuration      time.Duration
	maxRequestsToken   int
	blockDurationToken time.Duration
}

// Inicializador do estrategy redis
func NewRedisStrategy(repo *repository.RedisRepository, maxRequests int, blockDuration time.Duration, maxRequestsToken int, blockDurationToken time.Duration) Strategy {
	return &redisStrategy{
		repo:               repo,
		maxRequests:        maxRequests,
		blockDuration:      blockDuration,
		maxRequestsToken:   maxRequestsToken,
		blockDurationToken: blockDurationToken,
	}
}

// Envia para o repository as variáveis de acordo com a existência ou não do token na requisição
func (r *redisStrategy) Allow(ip string, token string) bool {
	if token != "" {
		return r.repo.Allow(ip, token, r.maxRequestsToken, r.blockDurationToken)
	}
	return r.repo.Allow(ip, "", r.maxRequests, r.blockDuration)
}
