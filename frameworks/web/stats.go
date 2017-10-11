package web

import (
	"encoding/json"
	"net/http"
)

func (r *RestAPI) HandleStatistics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := r.Statistics.Data()

	b, _ := json.Marshal(stats)

	w.Write(b)

}
