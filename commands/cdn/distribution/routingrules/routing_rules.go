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
			Use:     "routingrules",
			Aliases: []string{"rr"},
			Short:   "View the routing rules of a CDN distribution",
			Long: `Routing rules are the heart of a CDN distribution: each rule matches incoming requests by URL path prefix and scheme, then forwards them to an upstream origin, deciding per-rule whether to cache responses at the edge, run the Web Application Firewall, apply a per-IP rate-limit class, and allow/block traffic by country.

These sub-commands let you inspect a distribution's rules. Rules are created and replaced as a whole through 'ionosctl cdn ds create'/'update --routing-rules' (there is no add/remove-single-rule verb); use 'ionosctl cdn ds create --routing-rules-example' to see the JSON format.

The default columns are Scheme, Prefix, upstream Host, RateLimitClass and SniMode; add Caching, Waf, GeoRestrictionsAllowList or GeoRestrictionsBlockList with --cols to see the rest.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(GetDistributionRoutingRules())
	return cmd
}

func GetDistributionRoutingRules() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "routingrules",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "List the routing rules of a CDN distribution",
		LongDesc: `List the routing rules configured on a CDN distribution. Each row is one rule: the path prefix and scheme it matches, the upstream origin host it forwards to, and the rate-limit class and SNI mode in effect. Add --cols to reveal caching, WAF, and geo-restriction lists.

Use -o json to get the exact rule array, which you can edit and feed back to 'ionosctl cdn ds update --routing-rules' to change the rules.`,
		Example: "ionosctl cdn ds rr get --distribution-id ID",
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

			return c.Printer(allCols).Print(r.Properties.RoutingRules)
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagCDNDistributionID, constants.FlagIdShort, "", "The ID of the distribution whose routing rules to list",
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
