package main

import (
	"flag"
	"fmt"

	"github.com/riomhaire/lightauth2/usecases"
)

func main() {
	// Command lines
	username := flag.String("user", "anonymous", "Username associated with the token")
	password := flag.String("password", "", "Password to use - cannot be empty")
	roles := flag.String("roles", "guest:public", "List of roles separated by ':'")
	claim1 := flag.String("claim1", "", "Claim such as access token, fingerprint hash")
	claim2 := flag.String("claim2", "", "Claim such as access token, fingerprint hash")

	flag.Parse()
	if len(*password) == 0 {
		fmt.Println("Password Parameter cannot be empty")
		return
	}
	hash := usecases.HashPassword(*password, fmt.Sprintf("%v%v", *username, *password))
	fmt.Printf("%v,%v,true,%v,%v,%v\n", *username, hash, *roles, *claim1, *claim2)
}
