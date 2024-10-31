package rate_limiter

import (
	"time"

	"github.com/wandermaia/desafio-rate-limiter/internal/repository"
)

type redisStrategy struct {
	repo               *repository.RedisRepository
	maxRequests        int
	blockDuration      time.Duration
	maxRequestsToken   int
	blockDurationToken time.Duration
}

func NewRedisStrategy(repo *repository.RedisRepository, maxRequests int, blockDuration time.Duration, maxRequestsToken int, blockDurationToken time.Duration) Strategy {
	return &redisStrategy{
		repo:               repo,
		maxRequests:        maxRequests,
		blockDuration:      blockDuration,
		maxRequestsToken:   maxRequestsToken,
		blockDurationToken: blockDurationToken,
	}
}

func (r *redisStrategy) Allow(ip string, token string) bool {
	if token != "" {
		return r.repo.Allow(ip, token, r.maxRequestsToken, r.blockDurationToken)
	}
	return r.repo.Allow(ip, "", r.maxRequests, r.blockDuration)
}
