package interfaces

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"
)

type TestLogger struct {
	Buffer bytes.Buffer
}

// Append to String - so can check if required
func (d TestLogger) Log(level, message string) {
	d.Buffer.WriteString(fmt.Sprintf("[%s] %s\n", level, message))
}

func TestSuccessfulAuthenticate(t *testing.T) {
	registry := createTestRegistry()

	token, err := registry.AuthenticateInteractor.Authenticate("test", "secret")
	if err != nil || len(token) == 0 {
		t.Fail()
	}

}
func TestFailAuthenticate(t *testing.T) {
	registry := createTestRegistry()

	_, err := registry.AuthenticateInteractor.Authenticate("test", "bad")
	if err == nil {
		t.Fail()
	}

}

func TestFailAuthenticateUnknownUser(t *testing.T) {
	registry := createTestRegistry()

	_, err := registry.AuthenticateInteractor.Authenticate("unknown", "bad")
	if err == nil {
		t.Fail()
	}
	if err.Error() != "Unknown user" {
		t.Fail()
	}

}

func TestFailAuthenticateClaimsUnknownClaim(t *testing.T) {
	registry := createTestRegistry()

	claims := []string{"known"}
	_, err := registry.AuthenticateInteractor.AuthenticateClaims("test", claims)
	if err == nil {
		t.Fail()
	}

}

func TestFailAuthenticateClaimsSingleClaim(t *testing.T) {
	registry := createTestRegistry()

	claims := []string{"claim1"}
	token, err := registry.AuthenticateInteractor.AuthenticateClaims("test", claims)
	if err != nil || len(token) == 0 {
		t.Fail()
	}

}

func TestFailAuthenticateClaimsTwoClaim(t *testing.T) {
	registry := createTestRegistry()

	claims := []string{"claim1", "claim2"}
	token, err := registry.AuthenticateInteractor.AuthenticateClaims("test", claims)
	if err != nil || len(token) == 0 {
		t.Fail()
	}

}

func TestFailAuthenticateClaimsSingleClaimNoUser(t *testing.T) {
	registry := createTestRegistry()

	claims := []string{"claim1"}
	_, err := registry.AuthenticateInteractor.AuthenticateClaims("notest", claims)
	if err == nil {
		t.Fail()
	}
	if err.Error() != "Unknown user" {
		t.Fail()
	}

}

func TestFailAuthenticateClaimsNoClaims(t *testing.T) {
	registry := createTestRegistry()

	claims := []string{}
	_, err := registry.AuthenticateInteractor.AuthenticateClaims("test", claims)
	if err == nil {
		t.Fail()
	}
	if err.Error() != "No Claims" {
		t.Fail()
	}
}

func TestFailAuthenticateClaimsUnknownClaim2(t *testing.T) {
	registry := createTestRegistry()

	claims := []string{"claim1", "unknown"}
	_, err := registry.AuthenticateInteractor.AuthenticateClaims("test", claims)
	if err == nil {
		t.Fail()
	}

}

func TestFailValidate(t *testing.T) {
	registry := createTestRegistry()

	_, err := registry.TokenInteractor.ValidateToken("this is not a valid token")
	if err == nil {
		t.Fail()
	}
	if !strings.Contains(err.Error(), "invalid") {
		t.Fail()
	}

}

func TestSuccessValidate(t *testing.T) {
	registry := createTestRegistry()

	// Create token
	token, err := registry.AuthenticateInteractor.Authenticate("test", "secret")
	if err != nil || len(token) == 0 {
		t.Fail()
	}

	_, err = registry.TokenInteractor.ValidateToken(token)
	if err != nil {
		t.Fail()
	}

}

func TestSuccessDecode(t *testing.T) {
	registry := createTestRegistry()

	// Create token
	token, err := registry.AuthenticateInteractor.Authenticate("test", "secret")
	if err != nil || len(token) == 0 {
		t.Fail()
	}

	tok, err := registry.TokenInteractor.DecodeToken(token)
	if err != nil {
		t.Fail()
	}
	/*
	   	Id      string   `json:"id"`
	   	User    string   `json:"user"`
	   	Expires int64    `json:"expires"`
	   	Roles   []string `json:"roles"`
	   }
	*/
	if len(tok.User) == 0 {
		t.Fail()
	}
	if len(tok.Roles) == 0 {
		t.Fail()
	}

}

func createTestRegistry() usecases.Registry {
	logger := TestLogger{}

	database := NewTestDatabaseInteractor()
	database.Create(entities.User{"test", "939c1f673b7f5f5c991b4d4160642e72880e783ba4d7b04da260392f855214a6", true, []string{"user"}, "claim1", "claim2"})
	database.Create(entities.User{"admin", "50b911deac5df04e0a79ef18b04b29b245b8f576dcb7e5cca5937eb2083438ba", true, []string{"admin"}, "claim1", "claim2"})

	configuration := usecases.Configuration{}
	configuration.TokenTimeout = 3600
	configuration.SigningSecret = "secret"

	registry := usecases.Registry{}
	registry.Configuration = configuration
	registry.Logger = logger
	registry.StorageInteractor = database
	registry.AuthenticateInteractor = DefaultAuthenticateInteractor{&registry}
	registry.TokenInteractor = DefaultTokenInteractor{&registry}

	return registry
}

// This is a test implementation for test purposes
type TestDatabaseInteractor struct {
	db map[string]entities.User
}

func NewTestDatabaseInteractor() TestDatabaseInteractor {
	d := TestDatabaseInteractor{}
	d.db = make(map[string]entities.User)
	return d
}

func NewPopulatedTestDatabaseInteractorr(users []entities.User) TestDatabaseInteractor {
	d := TestDatabaseInteractor{}
	d.db = make(map[string]entities.User)

	for _, u := range users {
		d.db[u.Username] = u
	}

	return d
}

func (db TestDatabaseInteractor) Lookup(username string) (entities.User, error) {
	if val, ok := db.db[username]; ok {
		return val, nil
	} else {
		return entities.User{}, errors.New("Unknown user")
	}
}

func (db TestDatabaseInteractor) Create(user entities.User) error {
	if _, ok := db.db[user.Username]; ok {
		return errors.New("User exists")
	}
	db.db[user.Username] = user
	return nil
}

func (d TestDatabaseInteractor) add(u entities.User) {
	d.db[u.Username] = u
}
