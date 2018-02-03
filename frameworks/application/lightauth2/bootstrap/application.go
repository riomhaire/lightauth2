package bootstrap

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/riomhaire/lightauth2/frameworks"
	"github.com/riomhaire/lightauth2/frameworks/web"
	"github.com/riomhaire/lightauth2/interfaces"
	"github.com/riomhaire/lightauth2/usecases"
	"github.com/urfave/negroni"
)

const VERSION = "LightAuth2 Version 1.5.1"

type Application struct {
	registry *usecases.Registry
	restAPI  *web.RestAPI
}

func (a *Application) Initialize() {
	logger := frameworks.ConsoleLogger{}

	logger.Log("INFO", "Initializing")
	// Create Configuration
	configuration := usecases.Configuration{}
	sessionSecret := flag.String("sessionSecret", "secret", "Master key which is used to generate system jwt")
	sessionPeriod := flag.Int("sessionPeriod", 3600, "How many seconds before sessions expires")
	userFile := flag.String("usersFile", "users.csv", "List of Users and salted/hashed password with their roles")
	useSSL := flag.Bool("useSSL", false, "If True Enable SSL Server support")
	enableProfiling := flag.Bool("profile", false, "Enable profiling endpoint")
	serverCert := flag.String("serverCert", "server.crt", "Server Cert File")
	serverKey := flag.String("serverKey", "server.key", "Server Key File")

	port := flag.Int("port", 3030, "Port to use")

	enableKafkaLogging := flag.Bool("kafkaLogging", false, "Enable logging to Kafka")
	enableKafkaMetrics := flag.Bool("kafkaMetrics", false, "Enable metrics to Kafka")
	kafkaHost := flag.String("kafkaHost", "localhost", "Where Kafka is running")
	kafkaPort := flag.Int("kafkaPort", 9092, "Port where Kafka is listening")
	kafkaLoggingTopic := flag.String("kafkaLoggingTopic", "lightauth-logging", "Logging topic")
	kafkaMetricsTopic := flag.String("kafkaMetricsTopic", "lightauth-metrics", "Metrics topic")

	flag.Parse()
	// Set in config
	configuration.Application = "Authentication"
	configuration.Version = VERSION
	configuration.SigningSecret = *sessionSecret
	configuration.TokenTimeout = *sessionPeriod
	configuration.Store = *userFile
	configuration.SSL = *useSSL
	configuration.Profiling = *enableProfiling
	configuration.SSLCertificate = *serverCert
	configuration.SSLKey = *serverKey
	configuration.Port = *port
	configuration.KafkaLogging = *enableKafkaLogging
	configuration.KafkaMetrics = *enableKafkaMetrics
	configuration.KafkaHost = *kafkaHost
	configuration.KafkaPort = *kafkaPort
	configuration.KafkaLoggingTopic = *kafkaLoggingTopic
	configuration.KafkaMetricsTopic = *kafkaMetricsTopic

	registry := usecases.Registry{}
	a.registry = &registry
	registry.Configuration = configuration
	if configuration.KafkaLogging {
		registry.Logger = frameworks.NewKafkaLogger(configuration.KafkaHost, configuration.KafkaPort, configuration.KafkaLoggingTopic)
		logger.Log("Configuration", fmt.Sprintf("Using Kafa for logging. Server [%v:%v] Topic [%v]", configuration.KafkaHost, configuration.KafkaPort, configuration.KafkaLoggingTopic))
	} else {
		registry.Logger = logger
	}
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
