package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/riomhaire/lightauth2/frameworks"
	"github.com/riomhaire/lightauth2/frameworks/serviceregistry/consulagent"
	"github.com/riomhaire/lightauth2/frameworks/serviceregistry/defaultserviceregistry"
	"github.com/riomhaire/lightauth2/frameworks/web"
	"github.com/riomhaire/lightauth2/interfaces"
	"github.com/riomhaire/lightauth2/usecases"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

const VERSION = "LightAuth2 Version 1.8.8"

type Application struct {
	registry *usecases.Registry
	restAPI  *web.RestAPI
}

func (a *Application) Initialize(cmd *cobra.Command, args []string) {
	logger := frameworks.NewConsoleLogger(cmd.Flag("loggingLevel").Value.String())

	logger.Log(usecases.Info, "Initializing")
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
	configuration.UserAPI, _ = strconv.ParseBool(cmd.Flag("useUserAPI").Value.String())
	configuration.UserAPIHost = cmd.Flag("userAPIHost").Value.String()
	configuration.UserAPIKey = cmd.Flag("userAPIKey").Value.String()
	configuration.LoggingLevel = cmd.Flag("loggingLevel").Value.String()
	hostname, _ := os.Hostname()
	configuration.Host = hostname
	configuration.Consul, _ = strconv.ParseBool(cmd.Flag("consul").Value.String())
	configuration.ConsulHost = cmd.Flag("consulHost").Value.String()
	configuration.CacheTimeToLive, _ = strconv.Atoi(cmd.Flag("cacheTTL").Value.String())

	registry := usecases.Registry{}
	a.registry = &registry
	registry.Configuration = configuration
	registry.Logger = logger
	var database usecases.StorageInteractor
	if configuration.UserAPI {
		database = frameworks.NewUserAPIInteractor(&registry)
	} else {
		// Default CSV
		database = frameworks.NewCSVReaderDatabaseInteractor(&registry)
	}

	registry.StorageInteractor = database
	registry.AuthenticateInteractor = interfaces.DefaultAuthenticateInteractor{&registry}
	registry.TokenInteractor = interfaces.DefaultTokenInteractor{&registry}
	// Do we need external registry
	if configuration.Consul {

		registry.ExternalServiceRegistry = consulagent.NewConsulServiceRegistry(&registry, "/api/v2/authentication", "/api/v2/authentication/health")

	} else {
		registry.ExternalServiceRegistry = defaultserviceregistry.NewDefaultServiceRegistry(&registry)
	}

	// Create API
	restAPI := web.NewRestAPI(&registry)
	a.restAPI = &restAPI

	mux := http.NewServeMux()
	negroni := negroni.Classic()
	restAPI.Negroni = negroni

	// Add handlers
	mux.HandleFunc("/api/v2/authentication", restAPI.HandleAuthenticate)
	mux.HandleFunc("/api/v2/authentication/session", restAPI.HandleValidate)
	mux.HandleFunc("/api/v2/authentication/session/decoder", restAPI.HandleTokenDecode)
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
	handler := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler(mux) // Add coors
	negroni.UseHandler(handler)

	// Stats runs across all instances
	// n.UseFunc(AddWorkerHeader)

}

func (a *Application) Run() {
	a.registry.Logger.Log(usecases.Info, fmt.Sprintf("Running %s", a.registry.Configuration.Version))
	a.registry.Logger.Log(usecases.Info, a.registry.Configuration.String())
	// Register with external service if required ... default does nothing
	a.registry.ExternalServiceRegistry.Register()
	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}

func (a *Application) Stop() {
	a.registry.Logger.Log(usecases.Info, "Shutting Down REST API")
	a.registry.ExternalServiceRegistry.Deregister()
}
