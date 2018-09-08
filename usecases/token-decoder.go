package usecases

import (
	"errors"

	"github.com/riomhaire/lightauth2/entities"

	jwt "github.com/dgrijalva/jwt-go"
)

// Decode
func DecodeToken(tokenString, secret string) (entities.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != "HS256" {
			return entities.Token{}, errors.New("Unsupported Method")
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return entities.Token{}, err
	}
	//log.Println(token)
	session := entities.Token{}
	session.Id = tokenString
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		session.User = "unknown"
		session.Expires = -1
		roles := make([]string, 0)

		// Get user (if exists)
		if val, ok := claims["sub"]; ok {
			session.User = val.(string)
		}

		// Get Expires (if exists)
		if _, ok := claims["exp"]; ok {
			session.Expires = int64(claims["exp"].(float64))
		}

		// Get Roles (if exists)
		if _, ok := claims["roles"]; ok {
			croles := claims["roles"].([]interface{})
			for _, v := range croles {
				roles = append(roles, v.(string))
			}
			session.Roles = roles
		}
	}
	return session, nil
}
