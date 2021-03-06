package consulagent

import (
	"fmt"

	"github.com/riomhaire/consul"
	"github.com/riomhaire/lightauth2/usecases"
)

type ConsulServiceRegistry struct {
	registry       *usecases.Registry
	baseEndpoint   string
	healthEndpoint string
	id             string

	consulClient *consul.ConsulClient // This registers this service with consul - may extract this into a separate use case

}

func NewConsulServiceRegistry(registry *usecases.Registry, baseEndpoint, healthEndpoint string) *ConsulServiceRegistry {
	r := ConsulServiceRegistry{}
	r.registry = registry
	r.baseEndpoint = baseEndpoint
	r.healthEndpoint = healthEndpoint

	return &r
}

func (a *ConsulServiceRegistry) Register() error {
	// Register with consol (if required)
	if a.registry.Configuration.Consul {
		id := fmt.Sprintf("%v-%v-%v", a.registry.Configuration.Application, a.registry.Configuration.Host, a.registry.Configuration.Port)
		a.registry.Configuration.ConsulId = id // remember id for other system
		a.id = id                              // This is our safe copy

		a.consulClient, _ = consul.NewConsulClient(a.registry.Configuration.ConsulHost)
		health := fmt.Sprintf("http://%v:%v%v", a.registry.Configuration.Host, a.registry.Configuration.Port, a.healthEndpoint)
		a.registry.Logger.Log(usecases.Info, fmt.Sprintf("Registering with Consul at %v with %v %v", a.registry.Configuration.ConsulHost, a.baseEndpoint, health))

		a.consulClient.PeriodicRegister(id, a.registry.Configuration.Application, a.registry.Configuration.Host, a.registry.Configuration.Port, a.baseEndpoint, health, 15)
	}
	return nil

}

/*

 */
func (a *ConsulServiceRegistry) Deregister() error {
	if a.registry.Configuration.Consul {
		a.registry.Logger.Log(usecases.Info, fmt.Sprintf("De Registering %v with Consul at %v with %v ", a.id, a.registry.Configuration.ConsulHost, a.baseEndpoint))
		a.consulClient.DeRegister(a.id)
	}
	return nil
}
