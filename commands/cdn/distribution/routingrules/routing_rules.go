package routingrules

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allCols = []table.Column{
	{Name: "Scheme", JSONPath: "scheme", Default: true},
	{Name: "Prefix", JSONPath: "prefix", Default: true},
	{Name: "Host", JSONPath: "upstream.host", Default: true},
	{Name: "Caching", JSONPath: "upstream.caching"},
	{Name: "Waf", JSONPath: "upstream.waf"},
	{Name: "RateLimitClass", JSONPath: "upstream.rateLimitClass", Default: true},
	{Name: "SniMode", JSONPath: "upstream.sniMode", Default: true},
	{Name: "GeoRestrictionsAllowList", JSONPath: "upstream.geoRestrictions.allowList"},
	{Name: "GeoRestrictionsBlockList", JSONPath: "upstream.geoRestrictions.blockList"},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "routingrules",
			Aliases:          []string{"rr"},
			Short:            "Commands related to distribution routing rules",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

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

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, r.Properties.RoutingRules, cols))
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagCDNDistributionID, constants.FlagIdShort, "", "The ID of the distribution",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.DistributionsProperty(func(r cdn.Distribution) string {
				return r.Id
			})
		}, constants.CDNApiRegionalURL, constants.CDNLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
