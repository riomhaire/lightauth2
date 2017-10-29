package interfaces

import (
	"errors"
	"fmt"

	"github.com/riomhaire/lightauth2/usecases"
)

type DefaultAuthenticateInteractor struct {
	Registry *usecases.Registry
}

func (u DefaultAuthenticateInteractor) Authenticate(username string, password string) (string, error) {
	user, err := u.Registry.StorageInteractor.Lookup(username)
	if err != nil {
		return "", err
	}
	// OK we should compare whats entered with there ... passwords are encrypted
	hashedPassword := usecases.HashPassword(password, fmt.Sprintf("%v%v", username, password)) // I know each user should have on salt which is not the user

	//
	err = user.PasswordMatch(hashedPassword)
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

// Match via list of claims
func (u DefaultAuthenticateInteractor) AuthenticateClaims(username string, claims []string) (string, error) {
	user, err := u.Registry.StorageInteractor.Lookup(username)
	if err != nil {
		return "", err
	}

	// Check Claims
	if len(claims) == 0 {
		return "", errors.New("No Claims")
	}

	// Match claim 1 against claim 1
	if len(claims) > 0 && claims[0] != user.Claim1 {
		return "", errors.New("Authentication Failure")
	}

	// Match claim 2 against claim 2
	if len(claims) > 1 && claims[1] != user.Claim2 {
		return "", errors.New("Authentication Failure")
	}

	// OK create a token - Should be from config
	token, err := u.Registry.TokenInteractor.CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
