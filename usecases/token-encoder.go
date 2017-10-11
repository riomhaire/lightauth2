package usecases

import (
	"log"
	"time"

	"github.com/riomhaire/lightauth2/entities"

	jwt "github.com/dgrijalva/jwt-go"
)

func EncodeToken(user entities.User, secondsToLive int, secret string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	exp := time.Now().Add(time.Second * time.Duration(secondsToLive)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Username,
		"exp":   exp,
		"jid":   NewUUID(),
		"roles": user.Roles,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return tokenString, nil

}
