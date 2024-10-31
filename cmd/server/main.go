package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wandermaia/desafio-rate-limiter/internal/config"
	"github.com/wandermaia/desafio-rate-limiter/internal/handler"
	"github.com/wandermaia/desafio-rate-limiter/internal/middleware"
	"github.com/wandermaia/desafio-rate-limiter/internal/rate_limiter"
	"github.com/wandermaia/desafio-rate-limiter/internal/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	repo := repository.NewRedisRepository(cfg.RedisAddress, cfg.RedisPassword)
	strategy := rate_limiter.NewRedisStrategy(repo, cfg.MaxRequestsPerSecond, time.Duration(cfg.BlockDuration)*time.Second,
		cfg.MaxRequestsPerSecondToken, time.Duration(cfg.BlockDurationToken)*time.Second)
	limiter := rate_limiter.NewRateLimiter(strategy)

	log.Printf("MaxRequestsPerSecond: %v BlockDuration: %v , MaxRequestsPerSecondToken: %v , BlockDurationToken: %v", cfg.MaxRequestsPerSecond, cfg.BlockDuration,
		cfg.MaxRequestsPerSecondToken, cfg.BlockDurationToken)

	r := gin.Default()

	// recebe como parâmetro rate_limiter.RateLimiter, que é uma interface.
	r.Use(middleware.RateLimiterMiddleware(limiter))

	r.GET("/test", handler.TestHandler)

	r.Run()
}
