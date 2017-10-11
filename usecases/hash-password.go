package usecases

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

func HashPassword(password, salt string) string {
	bpassword := []byte(password)
	bsalt := []byte(salt)
	v := pbkdf2.Key(bpassword, bsalt, 4096, sha256.Size, sha256.New)
	hash := fmt.Sprintf("%x", v) // I know each user should have on salt
	//log.Printf("%s -> %s\n", password, hash)
	return hash
}
