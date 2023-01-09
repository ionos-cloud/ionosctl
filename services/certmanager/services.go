package certmanager

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/services/certmanager/resources"
	"github.com/spf13/viper"
)

type Services struct {
	// Certificate Manager Resources Services
	Certs   func() resources.CertsService
	Context context.Context
}

// InitClient for Commands
func (c *Services) InitClient() (*resources.Client, error) {
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token), // Token support
		config.GetServerUrl(),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}

// InitServices for Commands
func (c *Services) InitServices(client *resources.Client) error {
	c.Certs = func() resources.CertsService { return resources.NewCertsService(client, c.Context) }
	return nil
}
