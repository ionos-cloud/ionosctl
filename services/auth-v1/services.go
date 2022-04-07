package auth_v1

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/services/auth-v1/resources"
	"github.com/spf13/viper"
)

type Services struct {
	// Auth Resources Services
	Tokens func() resources.TokensService

	// Context
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
	c.Tokens = func() resources.TokensService { return resources.NewTokenService(client, c.Context) }
	return nil
}
