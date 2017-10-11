package interfaces

import (
	"fmt"

	"github.com/riomhaire/lightauth2/usecases"
)

type DefaultAuthenticateInteractor struct {
	Registry *usecases.Registry
}

func (u DefaultAuthenticateInteractor) Authenticate(username string, claims []string) (string, error) {
	user, err := u.Registry.StorageInteractor.Lookup(username)
	if err != nil {
		return "", err
	}
	// OK we should compare whats entered with there ... passwords are encrypted
	signedClaims := make([]string, 1, 1)
	signedClaims[0] = usecases.HashPassword(claims[0], fmt.Sprintf("%v%v", username, claims[0])) // I know each user should have on salt which is not the user

	err = user.ClaimsMatch(signedClaims)
	if err != nil {
		return "", err
	}

	// OK create a token - Should be from config
	token, err := u.Registry.TokenInteractor.CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
