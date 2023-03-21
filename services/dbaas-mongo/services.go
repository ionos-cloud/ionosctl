package dbaas_mongo

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-mongo/resources"
)

type Services struct {
	Clusters    func() resources.ClustersService
	Templates   func() resources.TemplatesService
	Users       func() resources.UsersService
	ApiMetadata func() resources.ApiMetadataService
	// Context
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices(client *client.Client) error {
	c.Clusters = func() resources.ClustersService { return resources.NewClustersService(client, c.Context) }
	c.Templates = func() resources.TemplatesService { return resources.NewTemplatesService(client, c.Context) }
	c.Users = func() resources.UsersService { return resources.NewUsersService(client, c.Context) }
	c.ApiMetadata = func() resources.ApiMetadataService { return resources.NewApiMetadataService(client, c.Context) }
	return nil
}
