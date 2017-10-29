package entities

import "errors"

type User struct {
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Enabled  bool     `json:"enabled,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Claim1   string   `json:"claim1,omitempty"`
	Claim2   string   `json:"claim2,omitempty"`
}

// PasswordMatch - in this implementations assume 1st claim is the password so we compare on that
func (u *User) PasswordMatch(password string) error {
	// Empty match
	if len(u.Password) == 0 && len(password) == 0 {
		// Valid
		return nil
	}

	// Non empty match
	if len(u.Password) > 0 && len(password) > 0 && u.Password == password {
		// Match
		return nil
	}

	return errors.New("Passwords do not match")
}

// ClaimsMatch - in this implementations assume 1st claim is the password so we compare on that
func (u *User) ClaimsMatch(claims []string) error {
	// Empty match
	if len(u.Password) == 0 && len(claims) == 0 {
		// Valid
		return nil
	}

	// Non empty match
	if len(u.Password) > 0 && len(claims) > 0 && u.Password == claims[0] {
		// Match
		return nil
	}

	return errors.New("Claims do not match")
}
