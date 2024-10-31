package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wandermaia/desafio-rate-limiter/internal/rate_limiter"
)

func RateLimiterMiddleware(limiter rate_limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		if !limiter.Allow(ip, token) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			c.Abort()
			return
		}

		c.Next()
	}
}
