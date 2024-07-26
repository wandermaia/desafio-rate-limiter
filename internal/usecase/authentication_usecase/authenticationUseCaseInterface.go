package authentication_usecase

// Interface para o caso de uso de autenticação
type AuthenticationUseCaseInterface interface {
	CreateToken(userid uint64) (*TokenDetailsDTO, error)
	CreateAuth(userid uint64, td *TokenDetailsDTO) error
	GetAuthByAccessUuid(accessUuid string) (string, error)
	DeleteAuth(givenUuid string) error
}
