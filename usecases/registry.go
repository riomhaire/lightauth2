package usecases

import (
	"fmt"

	"github.com/riomhaire/lightauth2/frameworks/serviceregistry"
)

// Configuration containing data from the environment which is used to define program behaviour
type Configuration struct {
	Application       string
	Version           string
	SigningSecret     string
	TokenTimeout      int
	Store             string
	SSL               bool
	Profiling         bool
	SSLCertificate    string
	SSLKey            string
	Port              int
	KafkaLogging      bool
	KafkaMetrics      bool
	KafkaHost         string
	KafkaPort         int
	KafkaLoggingTopic string
	KafkaMetricsTopic string
	UserAPI           bool
	UserAPIKey        string
	UserAPIHost       string
	LoggingLevel      string
	Host              string
	Consul            bool
	ConsulHost        string
	ConsulId          string // ID of this client
	CacheTimeToLive   int    // Seconds which to cache responses from the user API. 0 = dont cache
}

type Registry struct {
	Configuration           Configuration
	Logger                  Logger
	AuthenticateInteractor  AuthenticateInteractor
	StorageInteractor       StorageInteractor
	TokenInteractor         TokenInteractor
	ExternalServiceRegistry serviceregistry.ServiceRegistry
}

func (c *Configuration) String() string {
	return fmt.Sprintf("\nCONFIGURATION\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n",
		"Application",
		c.Application,
		"SigningSecret",
		c.SigningSecret,
		"TokenTimeout",
		c.TokenTimeout,
		"Store",
		c.Store,
		"SSL",
		c.SSL,
		"Port",
		c.Port,
		"UserAPI",
		c.UserAPI,
		"UserAPIHost",
		c.UserAPIHost,
		"LoggingLevel",
		c.LoggingLevel,
	)
}
