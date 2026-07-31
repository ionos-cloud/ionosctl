package distribution

import (
	"context"

	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func List() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "cdn",
			Resource:  "distribution",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "List CDN distributions across all locations",
			LongDesc: `List CDN distributions. By default all locations are queried and the results merged. Narrow the list with --domain (substring match on the served hostname) or --state (provisioning state).

The default columns show the distribution Id, Domain, bound CertificateId, and State (AVAILABLE, BUSY, FAILED, or UNKNOWN).`,
			Example: `# List every distribution
ionosctl cdn ds list

# List only distributions whose domain contains "example.com" and are ready to serve
ionosctl cdn ds list --domain example.com --state AVAILABLE`,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return nil
			},
			CmdRun: func(c *core.CommandConfig) error {
				return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
					cdnClient := cdn.NewAPIClient(cfg)
					req := cdnClient.DistributionsApi.DistributionsGet(context.Background())

					if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionFilterState); viper.IsSet(fn) {
						req = req.FilterState(viper.GetString(fn))
					}
					if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionFilterDomain); viper.IsSet(fn) {
						req = req.FilterDomain(viper.GetString(fn))
					}

					ls, _, err := req.Execute()
					return ls, err
				})
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagCDNDistributionFilterDomain, "", "", "Return only distributions whose served domain contains this value (substring match)")
	cmd.AddSetFlag(constants.FlagCDNDistributionFilterState, "", "", []string{"AVAILABLE", "BUSY", "FAILED", "UNKNOWN"}, "Return only distributions in this provisioning state")

	return cmd
}
