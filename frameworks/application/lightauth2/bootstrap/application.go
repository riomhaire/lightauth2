package bootstrap

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/riomhaire/lightauth2/frameworks"
	"github.com/riomhaire/lightauth2/frameworks/web"
	"github.com/riomhaire/lightauth2/interfaces"
	"github.com/riomhaire/lightauth2/usecases"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

const VERSION = "LightAuth2 Version 1.6"

type Application struct {
	registry *usecases.Registry
	restAPI  *web.RestAPI
}

func (a *Application) Initialize(cmd *cobra.Command, args []string) {
	logger := frameworks.ConsoleLogger{}

	logger.Log("INFO", "Initializing")
	// Create Configuration
	configuration := usecases.Configuration{}

	// Set in config
	configuration.Application = "Authentication"
	configuration.Version = VERSION
	configuration.SigningSecret = cmd.Flag("sessionSecret").Value.String()
	configuration.TokenTimeout, _ = strconv.Atoi(cmd.Flag("sessionPeriod").Value.String())
	configuration.Store = cmd.Flag("usersFile").Value.String()
	configuration.SSL, _ = strconv.ParseBool(cmd.Flag("useSSL").Value.String())
	configuration.Profiling, _ = strconv.ParseBool(cmd.Flag("profile").Value.String())
	configuration.SSLCertificate = cmd.Flag("serverCert").Value.String()
	configuration.SSLKey = cmd.Flag("serverKey").Value.String()
	configuration.Port, _ = strconv.Atoi(cmd.Flag("port").Value.String())
	registry := usecases.Registry{}
	a.registry = &registry
	registry.Configuration = configuration
	registry.Logger = logger
	database := frameworks.NewCSVReaderDatabaseInteractor(&registry)

	registry.StorageInteractor = database
	registry.AuthenticateInteractor = interfaces.DefaultAuthenticateInteractor{&registry}
	registry.TokenInteractor = interfaces.DefaultTokenInteractor{&registry}

	// Create API
	restAPI := web.NewRestAPI(&registry)
	a.restAPI = &restAPI

	mux := http.NewServeMux()
	negroni := negroni.Classic()
	restAPI.Negroni = negroni

	// Add handlers
	mux.HandleFunc("/api/v2/authentication", restAPI.HandleAuthenticate)
	mux.HandleFunc("/api/v2/session", restAPI.HandleValidate)
	mux.HandleFunc("/api/v2/session/decoder", restAPI.HandleTokenDecode)
	mux.HandleFunc("/api/v2/authentication/metrics", restAPI.HandleStatistics)
	mux.HandleFunc("/metrics", restAPI.HandleStatistics)
	mux.HandleFunc("/api/v2/authentication/health", restAPI.HandleHealth)
	mux.HandleFunc("/health", restAPI.HandleHealth)

	// Add Middleware
	if configuration.KafkaMetrics {
		negroni.UseFunc(restAPI.KafkaRecorder) // Record call in kafka
	}

	negroni.Use(restAPI.Statistics)
	negroni.UseFunc(restAPI.RecordCall)       // Calculates per second/minute rates
	negroni.UseFunc(restAPI.AddWorkerHeader)  // Add which instance
	negroni.UseFunc(restAPI.AddWorkerVersion) // Which version
	negroni.UseFunc(restAPI.AddCoorsHeader)   // Add coors
	negroni.UseHandler(mux)

	// Stats runs across all instances
	// n.UseFunc(AddWorkerHeader)

}

func (a *Application) Run() {
	a.registry.Logger.Log("INFO", fmt.Sprintf("Running %s", a.registry.Configuration.Version))
	a.registry.Logger.Log("INFO", a.registry.Configuration.String())
	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}
