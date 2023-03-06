package certmanager

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/services/certmanager/resources"
)

type Services struct {
	// Certificate Manager Resources Services
	Certs   func() resources.CertsService
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices(client *config.Client) error {
	c.Certs = func() resources.CertsService { return resources.NewCertsService(client, c.Context) }
	return nil
}
