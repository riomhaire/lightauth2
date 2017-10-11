package interfaces

import (
	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"
)

type DefaultTokenInteractor struct {
	Registry *usecases.Registry
}

func (s DefaultTokenInteractor) CreateToken(user entities.User) (string, error) {
	return usecases.EncodeToken(user, s.Registry.Configuration.TokenTimeout, s.Registry.Configuration.SigningSecret)

}

// Decode
func (s DefaultTokenInteractor) DecodeToken(token string) (entities.Token, error) {
	return usecases.DecodeToken(token, s.Registry.Configuration.SigningSecret)
}

// Validate
func (s DefaultTokenInteractor) ValidateToken(token string) (bool, error) {
	_, err := usecases.DecodeToken(token, s.Registry.Configuration.SigningSecret)
	if err != nil {
		return false, err
	}
	return true, nil
}
