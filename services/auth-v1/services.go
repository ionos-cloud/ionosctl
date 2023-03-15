package auth_v1

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
)

type Services struct {
	// Auth Resources Services
	Tokens func() resources.TokensService

	// Context
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices(client *client.Client) error {
	c.Tokens = func() resources.TokensService { return resources.NewTokenService(client, c.Context) }
	return nil
}
