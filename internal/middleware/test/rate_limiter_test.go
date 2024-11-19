package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wandermaia/desafio-rate-limiter/internal/middleware"
	"github.com/wandermaia/desafio-rate-limiter/internal/rate_limiter"
	"github.com/wandermaia/desafio-rate-limiter/internal/repository"
)

// Gera a configuração do router para os testes
func setupRouter() *gin.Engine {

	// Criação do repo, stratagey e limiter
	repo := repository.NewRedisRepository(viper.GetString("REDIS_ADDRESS"), viper.GetString("REDIS_PASSWORD"))
	strategy := rate_limiter.NewRedisStrategy(repo, viper.GetInt("MAX_REQUESTS"), time.Duration(viper.GetInt("BLOCK_DURATION"))*time.Second,
		viper.GetInt("MAX_REQUESTS_TOKEN"), time.Duration(viper.GetInt("BLOCK_DURATION_TOKEN"))*time.Second)
	limiter := rate_limiter.NewRateLimiter(strategy)

	//Limpando os dados do redis antes de inicializar o teste
	repo.FlushRedis()

	r := gin.Default()
	r.GET("/test", middleware.RateLimiterMiddleware(limiter), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "request successful",
		})
	})

	return r
}

// Função para a realização dos testes automatizados
func TestRateLimiterMiddleware(t *testing.T) {

	// Configurando as variáveis de ambiente para a execução dos testes
	viper.Set("MAX_REQUESTS", 2)
	viper.Set("BLOCK_DURATION", 60)
	viper.Set("MAX_REQUESTS_TOKEN", 5)
	viper.Set("BLOCK_DURATION_TOKEN", 60)
	viper.Set("REDIS_ADDRESS", "localhost:6379")
	viper.Set("REDIS_PASSWORD", "redis123")

	// Carregando as configurações do router
	router := setupRouter()

	t.Run("Requests abaixo do limite (IP)", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Real-IP", "127.0.0.1")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Esperado status 200, mas veio %d", w.Code)
			}
		}
	})

	t.Run("Requests acima do limite (IP)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Real-IP", "127.0.0.1")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Esperado status 429, mas veio %d", w.Code)
		}
	})

	t.Run("Requests abaixo do limite (Token)", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("API_KEY", "abc123")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Esperado status 200, mas veio %d", w.Code)
			}
		}
	})

	t.Run("Requests acima do limite (Token)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("API_KEY", "abc123")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Esperado status 429, mas veio %d", w.Code)
		}
	})
}
