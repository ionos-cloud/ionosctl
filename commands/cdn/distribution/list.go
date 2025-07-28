package distribution

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "The first element (of the total list of elements) to include in the response. Use together with limit for pagination")
	cmd.AddInt32Flag(constants.FlagMaxResults, "", 0, constants.DescMaxResults)

	cmd.Command.PersistentFlags().StringSlice(
		constants.ArgCols, nil,
		fmt.Sprintf(
			"Set of columns to be printed on output \nAvailable columns: %v",
			allCols,
		),
	)

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
			if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
				req = req.Offset(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(fn) {
				req = req.Limit(viper.GetInt32(fn))
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

	convertedItems, err := json2table.ConvertJSONToTable("", jsonpaths.CDNDistribution, items)
	if err != nil {
		return fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	out, err := jsontabwriter.GenerateOutputPreconverted(ls, convertedItems, tabheaders.GetHeaders(allCols, defaultCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
