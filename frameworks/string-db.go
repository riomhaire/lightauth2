package frameworks

import (
	"errors"

	"github.com/riomhaire/lightauth2/entities"
)

// This is a test implementation for test purposes
type StringDatabaseInteractor struct {
	db map[string]entities.User
}

func NewStringDatabaseInteractor() StringDatabaseInteractor {
	d := StringDatabaseInteractor{}
	d.db = make(map[string]entities.User)
	return d
}

func NewPopulatedStringDatabaseInteractor(users []entities.User) StringDatabaseInteractor {
	d := StringDatabaseInteractor{}
	d.db = make(map[string]entities.User)

	for _, u := range users {
		d.db[u.Username] = u
	}

	return d
}

func (db StringDatabaseInteractor) Lookup(username string) (entities.User, error) {
	if val, ok := db.db[username]; ok {
		return val, nil
	} else {
		return entities.User{}, errors.New("Unknown user")
	}
}

func (db StringDatabaseInteractor) Create(user entities.User) error {
	if _, ok := db.db[user.Username]; ok {
		return errors.New("User exists")
	}
	db.db[user.Username] = user
	return nil
}

func (d StringDatabaseInteractor) add(u entities.User) {
	d.db[u.Username] = u
}
