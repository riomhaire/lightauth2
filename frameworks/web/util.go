package web

import (
	"errors"
	"strings"
)

func extractAuthorization(header string) (string, error) {
	if len(header) < len(bearerPrefix) {
		return "", errors.New("Not Authorized")
	}

	prefix := strings.ToLower(header[:len(bearerPrefix)])

	if !strings.HasPrefix(prefix, bearerPrefix) {
		return "", errors.New("Not Authorized")
	}
	token := header[len(bearerPrefix):]
	return token, nil
}
