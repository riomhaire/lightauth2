package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/riomhaire/lightauth2/usecases"
)

type AuthenticateResponse struct {
	Token string `json:"token,omitempty"`
}

type LoginParameters struct {
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Claims   []string `json:"claims,omitempty"`
}

func (r *RestAPI) HandleAuthenticate(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	// Decode request
	var loginParameters LoginParameters
	_ = json.NewDecoder(req.Body).Decode(&loginParameters)
	// Call interactor - which one is dependent on whether password is present and claims
	var token string
	var err error
	if len(loginParameters.Password) > 0 { // authenticate via password
		token, err = r.Registry.AuthenticateInteractor.Authenticate(loginParameters.Username, loginParameters.Password)
	} else if len(loginParameters.Claims) > 0 { // authenticate via claims // Could be password - but not always
		token, err = r.Registry.AuthenticateInteractor.AuthenticateClaims(loginParameters.Username, loginParameters.Claims)
	} else {
		// Unknown
		err = errors.New("Not Authenticated")
	}

	// Decode result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		r.Registry.Logger.Log(usecases.Error, err.Error())
		w.WriteHeader(http.StatusUnauthorized) // unprocessable entity
		json.NewEncoder(w).Encode(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Authorization", fmt.Sprintf("%s %s", bearerPrefix, token))
		json.NewEncoder(w).Encode(AuthenticateResponse{token})
	}
}
