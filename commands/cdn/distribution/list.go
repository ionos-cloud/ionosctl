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
			ShortDesc: "Retrieve all distributions using pagination and optional filters",
			Example:   `ionosctl cdn ds list`,
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

	cmd.AddStringFlag(constants.FlagCDNDistributionFilterDomain, "", "", "Filter used to fetch only the records that contain specified domain.")
	cmd.AddSetFlag(constants.FlagCDNDistributionFilterState, "", "", []string{"AVAILABLE", "BUSY", "FAILED", "UNKNOWN"}, "Filter used to fetch only the records that contain specified state.")

	return cmd
}
