package routingrules

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allRoutingRulesColumns     = []string{"Scheme", "Prefix", "Host", "Caching", "Waf", "RateLimitClass", "SniMode", "GeoRestrictionsAllowList", "GeoRestrictionsBlockList"}
	defaultRoutingRulesColumns = []string{"Scheme", "Prefix", "Host", "RateLimitClass", "SniMode"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "routingrules",
			Aliases:          []string{"rr"},
			Short:            "Commands related to distribution routing rules",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(GetDistributionRoutingRules())
	return cmd
}

func GetDistributionRoutingRules() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "routingrules",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a distribution routing rules",
		Example:   "ionosctl cdn ds rr get --distribution-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagCDNDistributionID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			distributionID := viper.GetString(core.GetFlagName(c.NS, constants.FlagCDNDistributionID))
			r, _, err := client.Must().CDNClient.DistributionsApi.DistributionsFindById(context.Background(),
				distributionID).Execute()
			if err != nil {
				return err
			}

			if r.Properties.RoutingRules == nil {
				return nil
			}

			convertedItems, err := json2table.ConvertJSONToTable("", jsonpaths.CDNRoutingRule, *r.Properties.RoutingRules)
			if err != nil {
				return fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(*r.Properties.RoutingRules, convertedItems, tabheaders.GetHeaders(allRoutingRulesColumns, defaultRoutingRulesColumns, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagCDNDistributionID, constants.FlagIdShort, "", "The ID of the distribution",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.DistributionsProperty(func(r cdn.Distribution) string {
				return *r.Id
			})
		}, constants.CDNApiRegionalURL, constants.CDNLocations),
	)
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
