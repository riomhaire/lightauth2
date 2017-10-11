package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthenticateResponse struct {
	Token string `json:"token,omitempty"`
}

type LoginParameters struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *RestAPI) HandleAuthenticate(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	// Decode request
	var loginParameters LoginParameters
	_ = json.NewDecoder(req.Body).Decode(&loginParameters)
	// Call interactor
	token, err := r.Registry.AuthenticateInteractor.Authenticate(loginParameters.Username, loginParameters.claims())
	// Decode result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		r.Registry.Logger.Log("ERROR", err.Error())
		w.WriteHeader(http.StatusUnauthorized) // unprocessable entity
		json.NewEncoder(w).Encode(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Authorization", fmt.Sprintf("%s %s", bearerPrefix, token))
		json.NewEncoder(w).Encode(AuthenticateResponse{token})
	}
}

func (l *LoginParameters) claims() []string {
	var a = make([]string, 1)
	a[0] = (*l).Password
	return a
}
