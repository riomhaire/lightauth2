package web

import (
	"encoding/json"
	"net/http"
)

func (r *RestAPI) HandleValidate(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)

	auth, err := extractAuthorization(req.Header.Get("Authorization"))

	if err != nil || len(auth) == 0 {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	// Call interactor
	_, err = r.Registry.TokenInteractor.ValidateToken(auth)
	// Decode result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		r.Registry.Logger.Log("ERROR", err.Error())
		w.WriteHeader(http.StatusUnauthorized) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (r *RestAPI) HandleTokenDecode(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)

	auth, err := extractAuthorization(req.Header.Get("Authorization"))

	if err != nil || len(auth) == 0 {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	// Call interactor
	data, err := r.Registry.TokenInteractor.DecodeToken(auth)
	// Decode result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		r.Registry.Logger.Log("ERROR", err.Error())
		w.WriteHeader(http.StatusUnauthorized) // unprocessable entity
		json.NewEncoder(w).Encode(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
