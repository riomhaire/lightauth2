package main

import (
	"flag"
	"fmt"

	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"

	"strings"
	"time"
)

func main() {
	// Command lines
	username := flag.String("user", "anonymous", "Username associated with the token")
	roles := flag.String("roles", "guest:public", "List of roles separated by ':'")
	secret := flag.String("secret", "secret", "Key used to generate sessions")
	timeToLive := flag.Int("sessionPeriod", 3600, "How many seconds before sessions expires")
	tokenToDecode := flag.String("token", "", "If populated means decode token")

	flag.Parse()

	if len(*tokenToDecode) > 0 {
		// DECODE
		t, err := usecases.DecodeToken(*tokenToDecode, *secret)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("user    : %v\nexpires : %v\nroles   : %v\n", t.User, time.Unix(t.Expires, 0), t.Roles)
		}
	} else {
		// ENDCODE
		user := entities.User{}
		user.Username = *username
		r := strings.Split(*roles, ":")
		user.Roles = r

		token, err := usecases.EncodeToken(user, *timeToLive, *secret)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(token)
		}
	}

}
