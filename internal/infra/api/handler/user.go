package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wandermaia/desafio-rate-limiter/internal/usecase/authentication_usecase"
)

// Handler de albuns
type AccessHandler struct {
	authenticationUseCase authentication_usecase.AuthenticationUseCaseInterface
}

// Construtora do handler
func NewAccessHandler(authUseCaseInterface authentication_usecase.AuthenticationUseCaseInterface) *AccessHandler {
	return &AccessHandler{
		authenticationUseCase: authUseCaseInterface,
	}
}

// Struct para representar um usuário
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// Foi criado um usuário manualmente apenas para ser utilizando no token do rate limit
var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
	Phone:    "49123454322",
}

// Função para verificar o login
func (ah AccessHandler) Login(c *gin.Context) {

	// Usuário enviado na requisição para geração de token
	var userRequest User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		log.Print(err)
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	//Validando se as credenciais estão corretas
	if user.Username != userRequest.Username || user.Password != userRequest.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	//Invocando a função de criação de token
	token, err := ah.authenticationUseCase.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
