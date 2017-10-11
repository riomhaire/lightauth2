package usecases

import "github.com/riomhaire/lightauth2/entities"

// This file contains the various interface contracts used by the system.

type AuthenticateInteractor interface {
	Authenticate(user string, claims []string) (string, error)
}

type Logger interface {
	Log(level, message string)
}

type TokenInteractor interface {
	CreateToken(user entities.User) (string, error)
	DecodeToken(token string) (entities.Token, error)
	ValidateToken(token string) (bool, error)
}

type StorageInteractor interface {
	Lookup(username string) (entities.User, error)
	Create(user entities.User) error
}
