package auth_midlleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tokenutils "github.com/wandermaia/desafio-rate-limiter/pkg/tokenUtils"
)

// Midleware para incluir a autenticação nas rotas
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Extraindo o token
		token := tokenutils.ExtractToken(ctx.Request)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}

		// Validando o token
		err := tokenutils.TokenValid(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
