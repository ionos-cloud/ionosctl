package dbaas_postgres

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
)

type Services struct {
	// Dbaas Pgsql Resources Services
	Clusters func() resources.ClustersService
	// Context
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices(client *client.Client) error {
	c.Clusters = func() resources.ClustersService { return resources.NewClustersService(client, c.Context) }
	return nil
}
