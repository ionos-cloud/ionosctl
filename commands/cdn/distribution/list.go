package distribution

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/spf13/viper"
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
				return listDistributions(c)
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagCDNDistributionFilterDomain, "", "", "Filter used to fetch only the records that contain specified domain.")
	cmd.AddSetFlag(constants.FlagCDNDistributionFilterState, "", "", []string{"AVAILABLE", "BUSY", "FAILED", "UNKNOWN"}, "Filter used to fetch only the records that contain specified state.")

	return cmd
}

func listDistributions(c *core.CommandConfig) error {
	ls, err := completer.Distributions(
		func(req cdn.ApiDistributionsGetRequest) (cdn.ApiDistributionsGetRequest, error) {
			if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionFilterState); viper.IsSet(fn) {
				req = req.FilterState(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionFilterDomain); viper.IsSet(fn) {
				req = req.FilterDomain(viper.GetString(fn))
			}
			return req, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed listing cdn distributions: %w", err)
	}

	items, ok := ls.GetItemsOk()
	if !ok || items == nil {
		return fmt.Errorf("could not retrieve distributions")
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, items, cols))
}
