package token

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

var allTokenCols = []table.Column{
	{Name: "Token", JSONPath: "token", Default: true},
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)
}

func RunServerTokenGet(c *core.CommandConfig) error {
	c.Verbose("ServerToken with id: %v is getting...", c.Flags().String(cloudapiv6.ArgServerId))

	t, _, err := c.CloudApiV6Services.Servers().GetToken(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if err != nil {
		return err
	}

	return c.Printer(allTokenCols).Print(t.Token)
}
