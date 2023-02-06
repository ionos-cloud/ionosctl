package container_registry

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/services/container-registry/resources"
)

type Services struct {
	// Container Registry Resources Services
	Registry func() resources.RegistriesService
	// Context
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices(client *config.Client) error {
	c.Registry = func() resources.RegistriesService {
		return resources.NewRegistriesService(client, c.Context)
	}
	return nil
}
