package web

import (
	"net/http"
	"os"
)

// AddWorkerHeader - adds header of which node actually processed request
func (r *RestAPI) AddWorkerHeader(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}
	rw.Header().Add("X-Worker", host)
	next(rw, req)
}

// AddWorkerVersion - adds header of which version is installed
func (r *RestAPI) AddWorkerVersion(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	rw.Header().Add("X-Worker-Version", r.Registry.Configuration.Version)
	next(rw, req)
}
