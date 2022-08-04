package dataplatform

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"
	"github.com/spf13/viper"
)

type Services struct {
	// Data Platform Resources Services
	Clusters func() resources.ClustersService
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
	c.Clusters = func() resources.ClustersService { return resources.NewClustersService(client, c.Context) }
	return nil
}
