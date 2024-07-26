package authentication_usecase

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"github.com/wandermaia/desafio-rate-limiter/internal/infra/cache"
)

// Definições dos tokens
type TokenDetailsDTO struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

const (
	AlbumPrefixKeyCache = "token_"
	CacheTTL            = "TOKEN_TTL_SECONDS" // Nome da chave da variável de ambiente
)

// UseCase para a autenticação. Utiliza uma interface do cache para manter o desacoplamento
type AuthenticationUseCase struct {
	authenticationCacheInterface cache.CacheInterface
}

// Função "Construtora"
func NewAuthenticationUseCase(cacheInterface cache.CacheInterface) *AuthenticationUseCase {
	return &AuthenticationUseCase{
		authenticationCacheInterface: cacheInterface,
	}
}

// Gera um novo token
func (auc AuthenticationUseCase) CreateToken(userid uint64) (*TokenDetailsDTO, error) {

	// Definindo o id e os tempos do token
	tokenDetails := &TokenDetailsDTO{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.AccessUuid = uuid.NewV4().String()

	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetails.RefreshUuid = uuid.NewV4().String()

	var err error

	//Cria o token utilizando a chave
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["user_id"] = userid
	accessTokenClaims["access_uuid"] = tokenDetails.AccessUuid
	accessTokenClaims["exp"] = tokenDetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// Gerando a assinatura do access token
	tokenDetails.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Cria o Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokenDetails.RtExpires

	// Gerando a assinatura do refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetails.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	// Retornando todos os detalhes do token
	return tokenDetails, nil
}

// Função para criar a autenticação do token
func (auc AuthenticationUseCase) CreateAuth(userid uint64, tokenDetailsDTO *TokenDetailsDTO) error {

	// Convertendo os tempos de expiração para UTC
	accesstoken := time.Unix(tokenDetailsDTO.AtExpires, 0)
	refreshToken := time.Unix(tokenDetailsDTO.RtExpires, 0)
	now := time.Now()

	// Inserindo os dados no cache
	err := auc.authenticationCacheInterface.Set(tokenDetailsDTO.AccessUuid, strconv.Itoa(int(userid)), accesstoken.Sub(now))
	if err != nil {
		return err
	}
	errRefresh := auc.authenticationCacheInterface.Set(tokenDetailsDTO.RefreshUuid, strconv.Itoa(int(userid)), refreshToken.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (auc AuthenticationUseCase) GetAuthByAccessUuid(accessUuid string) (string, error) {

	userId, err := auc.authenticationCacheInterface.Get(accessUuid)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (auc AuthenticationUseCase) DeleteAuth(givenUuid string) error {
	err := auc.authenticationCacheInterface.Delete(givenUuid)
	if err != nil {
		return err
	}
	return nil
}
