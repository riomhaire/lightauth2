package web

import (
	"github.com/riomhaire/lightauth2/usecases"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

var bearerPrefix = "bearer "

type RestAPI struct {
	Registry   *usecases.Registry
	Statistics *stats.Stats
	Negroni    *negroni.Negroni
}

func NewRestAPI(registry *usecases.Registry) RestAPI {
	api := RestAPI{}
	api.Registry = registry
	api.Statistics = stats.New()

	return api
}
