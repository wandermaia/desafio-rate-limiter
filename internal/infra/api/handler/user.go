package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wandermaia/desafio-rate-limiter/internal/usecase/authentication_usecase"
	tokenutils "github.com/wandermaia/desafio-rate-limiter/pkg/tokenUtils"
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

// Armazena os metadados (access_uuid e user_id) que precisaremos fazer uma busca no redis
type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

// Função para autenticar o usuário
func (ah AccessHandler) CreateAccessToken(c *gin.Context) {

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
	tokenDetailsDTO, err := ah.authenticationUseCase.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = ah.authenticationUseCase.CreateAuth(userRequest.ID, tokenDetailsDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Segregando apenas os tokens para retornar ao usuário
	tokens := map[string]string{
		"access_token":  tokenDetailsDTO.AccessToken,
		"refresh_token": tokenDetailsDTO.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

// Função para logout
func (ah AccessHandler) Logout(ctx *gin.Context) {

	// Validando o token
	token := tokenutils.ExtractToken(ctx.Request)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// Extraindo os metadados do token
	tokenAuth, err := ExtractTokenMetadata(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// Deletando a autenticação
	delErr := ah.authenticationUseCase.DeleteAuth(tokenAuth.AccessUuid)
	if delErr != nil {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx.JSON(http.StatusOK, "Successfully logged out")
}

// Função para teste
func (ah AccessHandler) Health(ctx *gin.Context) {

	// // Extraindo o token da request
	// token := tokenutils.ExtractToken(ctx.Request)
	// if token == "" {
	// 	ctx.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	// // Extraindo os metadados do token
	// tokenAuth, err := ExtractTokenMetadata(token)
	// if err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	// // Verificando a autenticação
	// userId, err := ah.FetchAuth(tokenAuth)
	// if err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	// Segregando apenas os tokens para retornar ao usuário
	health := map[string]string{
		"status": "health",
		// "user_id": string(userId),
	}

	ctx.JSON(http.StatusOK, health)
}

// Extrai os metadados do token
func ExtractTokenMetadata(tokenString string) (*AccessDetails, error) {

	// Verificando se o token é válido e se nãoe stá expirado
	token, err := tokenutils.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

// Realiza a busca de metadados no redis
// FetchAuth() aceita os AccessDetails da função ExtractTokenMetadata, depois procura no redis.
// Se o registro não for encontrado, isso pode significar que o token expirou, portanto um erro é atirado.
func (ah AccessHandler) FetchAuth(authD *AccessDetails) (uint64, error) {
	//userid, err := client.Get(authD.AccessUuid).Result()
	userid, err := ah.authenticationUseCase.GetAuthByAccessUuid(authD.AccessUuid)
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

// Realiza o refresh do token
func (ah AccessHandler) RefreshAccessToken(ctx *gin.Context) {

	// Recuperando so tokens a partir do body
	mapToken := map[string]string{}
	if err := ctx.ShouldBindJSON(&mapToken); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//Validando o token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Confirma a assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if token.Claims.Valid() != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token

		// Deletando a autenticação
		delErr := ah.authenticationUseCase.DeleteAuth(refreshUuid)
		if delErr != nil {
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		//Invocando a função de criação de token
		tokenDetailsDTO, err := ah.authenticationUseCase.CreateToken(userId)
		if err != nil {
			ctx.JSON(http.StatusForbidden, err.Error())
			return
		}

		//save the tokens metadata to redis
		err = ah.authenticationUseCase.CreateAuth(userId, tokenDetailsDTO)
		if err != nil {
			ctx.JSON(http.StatusForbidden, err.Error())
			return
		}

		tokens := map[string]string{
			"access_token":  tokenDetailsDTO.AccessToken,
			"refresh_token": tokenDetailsDTO.RefreshToken,
		}
		ctx.JSON(http.StatusCreated, tokens)
	} else {
		ctx.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
