package usecases

import "fmt"

// Configuration containing data from the environment which is used to define program behaviour
type Configuration struct {
	Version        string
	SigningSecret  string
	TokenTimeout   int
	Store          string
	SSL            bool
	Profiling      bool
	SSLCertificate string
	SSLKey         string
	Port           int
}

type Registry struct {
	Configuration          Configuration
	Logger                 Logger
	AuthenticateInteractor AuthenticateInteractor
	StorageInteractor      StorageInteractor
	TokenInteractor        TokenInteractor
}

func (c *Configuration) String() string {
	return fmt.Sprintf("\nCONFIGURATION\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n",
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
	)
}
