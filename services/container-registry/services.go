package container_registry

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
)

type Services struct {
	// Container Registry Resources Services
	Registry   func() resources.RegistriesService
	Token      func() resources.TokenService
	Location   func() resources.LocationsService
	Name       func() resources.NameService
	Repository func() resources.RepositoryService
	// Context
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices(client *config.Client) error {
	c.Registry = func() resources.RegistriesService {
		return resources.NewRegistriesService(client, c.Context)
	}
	c.Token = func() resources.TokenService {
		return resources.NewTokenService(client, c.Context)
	}
	c.Location = func() resources.LocationsService {
		return resources.NewLocationsService(client, c.Context)
	}
	c.Name = func() resources.NameService {
		return resources.NewNameService(client, c.Context)
	}
	c.Repository = func() resources.RepositoryService {
		return resources.NewRepositoryService(client, c.Context)
	}
	return nil
}
