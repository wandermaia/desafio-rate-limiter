package tokenutils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Verifica token
func VerifyToken(tokenString string) (*jwt.Token, error) {

	// Verificando se o token é válido
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validando a assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Função para extrair o token a partir do header
func ExtractToken(r *http.Request) string {

	// Coletando o header que contem a chave Authorization e segregando a parte do token
	// bearToken := r.Header.Get("Authorization")
	// apiKey := r.Header.Get("API_KEY")
	// strArr := strings.Split(bearToken, " ")
	// if len(strArr) == 2 {
	// 	return strArr[1]
	// }
	// return ""
	return r.Header.Get("API_KEY")
}

// Validação do token
func TokenValid(tokenString string) error {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return err
	}

	if token.Claims.Valid() != nil || !token.Valid {
		return err
	}

	return nil
}
