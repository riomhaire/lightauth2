package frameworks

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"
)

// This is a test implementation for test purposes
type CSVReaderDatabaseInteractor struct {
	registry *usecases.Registry
	db       map[string]entities.User
}

func NewCSVReaderDatabaseInteractor(registry *usecases.Registry) CSVReaderDatabaseInteractor {
	d := CSVReaderDatabaseInteractor{}
	d.db = make(map[string]entities.User)
	d.registry = registry

	log.Printf("Reading User Database %s\n", registry.Configuration.Store)
	d.db, _ = loadUsers(registry.Configuration.Store)
	log.Printf("#Number of users = %v\n", len(d.db))
	return d
}

func (db CSVReaderDatabaseInteractor) Lookup(username string) (entities.User, error) {
	if val, ok := db.db[username]; ok {
		return val, nil
	} else {
		return entities.User{}, errors.New("Unknown user")
	}
}

func (db CSVReaderDatabaseInteractor) Create(user entities.User) error {
	if _, ok := db.db[user.Username]; ok {
		return errors.New("User exists")
	}
	db.db[user.Username] = user
	return nil
}

func (d CSVReaderDatabaseInteractor) add(u entities.User) {
	d.db[u.Username] = u
}

// Initiaizes data structues - IE Read user DB
func loadUsers(filename string) (map[string]entities.User, error) {
	users := make(map[string]entities.User)
	csvfile, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
		return users, err
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		return users, err
	}
	// Create user map
	for index, row := range records {
		if index > 0 && len(row) > 0 {
			user := entities.User{}
			user.Username = row[0]
			user.Password = row[1]

			v, _ := strconv.ParseBool(row[2])
			user.Enabled = v
			roles := strings.Split(row[3], ":")
			user.Roles = roles
			// Add
			users[user.Username] = user
		}
	}
	return users, nil
}
