package dbaas_pg

import (
	"context"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/services/dbaas-pg/resources"
	"github.com/spf13/viper"
)

type Services struct {
	// Dbaas Pgsql Resources Services
	Clusters func() resources.ClustersService
	Backups  func() resources.BackupsService
	Versions func() resources.VersionsService
	Infos    func() resources.InfosService
	Restores func() resources.RestoresService
	Logs     func() resources.LogsService
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
	c.Backups = func() resources.BackupsService { return resources.NewBackupsService(client, c.Context) }
	c.Versions = func() resources.VersionsService { return resources.NewVersionsService(client, c.Context) }
	c.Infos = func() resources.InfosService { return resources.NewInfosService(client, c.Context) }
	c.Restores = func() resources.RestoresService { return resources.NewRestoresService(client, c.Context) }
	c.Logs = func() resources.LogsService { return resources.NewLogsService(client, c.Context) }
	return nil
}
