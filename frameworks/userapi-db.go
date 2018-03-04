package frameworks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"
	"gopkg.in/resty.v1"
)

// This file contains the code to interact with the userAPI as a user source
// Rather than the default csv
type UserAPIInteractor struct {
	registry *usecases.Registry
}

func NewUserAPIInteractor(registry *usecases.Registry) UserAPIInteractor {
	d := UserAPIInteractor{}
	d.registry = registry

	log.Printf("Reading UserAPI at %s\n", registry.Configuration.UserAPIHost)
	return d
}

func (d UserAPIInteractor) Lookup(username string) (entities.User, error) {
	target := fmt.Sprintf("%v/api/v1/user/account/%v", d.registry.Configuration.UserAPIHost, username)
	token := d.registry.Configuration.UserAPIKey
	d.registry.Logger.Log(usecases.Trace, fmt.Sprintf("Requesting '%v' from '%v'", username, target))
	resp, err := resty.R().SetHeader("Accept", "application/json").SetAuthToken(token).Get(target)
	user := entities.User{}
	if err != nil {
		d.registry.Logger.Log(usecases.Error, err.Error())
		return user, err
	}
	// Decode
	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		d.registry.Logger.Log(usecases.Error, err.Error())
	}
	return user, err
}

// We dont use this
func (d UserAPIInteractor) Create(user entities.User) error {
	return errors.New("Not Needed/Implemented For UserAPI")
}
