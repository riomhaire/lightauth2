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
	if next != nil {
		next(rw, req)
	}
}

// AddWorkerVersion - adds header of which version is installed
func (r *RestAPI) AddWorkerVersion(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	version := r.Registry.Configuration.Version
	if len(version) == 0 {
		version = "UNKNOWN"
	}
	rw.Header().Add("X-Worker-Version", version)
	if next != nil {
		next(rw, req)
	}
}
