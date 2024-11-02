package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wandermaia/desafio-rate-limiter/internal/handler"
	"github.com/wandermaia/desafio-rate-limiter/internal/middleware"
	"github.com/wandermaia/desafio-rate-limiter/internal/rate_limiter"
	"github.com/wandermaia/desafio-rate-limiter/internal/repository"
)

func main() {

	// Carregando as variáveis de ambiente
	viper.AutomaticEnv()

	// viper.Set("MAX_REQUESTS", 5)
	// viper.Set("BLOCK_DURATION", 60)
	// viper.Set("MAX_REQUESTS_TOKEN", 10)
	// viper.Set("BLOCK_DURATION_TOKEN", 60)
	// viper.Set("REDIS_ADDRESS", "localhost:6379")
	// viper.Set("REDIS_PASSWORD", "redis123")

	// Criação do repo, stratagey e limiter
	repo := repository.NewRedisRepository(viper.GetString("REDIS_ADDRESS"), viper.GetString("REDIS_PASSWORD"))
	strategy := rate_limiter.NewRedisStrategy(repo, viper.GetInt("MAX_REQUESTS"), time.Duration(viper.GetInt("BLOCK_DURATION"))*time.Second,
		viper.GetInt("MAX_REQUESTS_TOKEN"), time.Duration(viper.GetInt("BLOCK_DURATION_TOKEN"))*time.Second)
	limiter := rate_limiter.NewRateLimiter(strategy)

	// Criação do server
	router := gin.Default()
	router.GET("/test", middleware.RateLimiterMiddleware(limiter), handler.TestHandler)
	router.Run(viper.GetString("PORT"))
}
