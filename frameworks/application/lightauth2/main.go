package main

import (
	"github.com/riomhaire/lightauth2/frameworks/application/lightauth2/bootstrap"
)

func main() {

	application := bootstrap.Application{}

	application.Initialize()
	application.Run()

}
