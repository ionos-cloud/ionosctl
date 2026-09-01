package distribution

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "distribution",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a distribution's domain, certificate binding, or routing rules",
		LongDesc: `Update an existing CDN distribution. Only the properties you pass are changed: the command first GETs the current distribution, overlays the flags you set (--domain, --certificate-id, --routing-rules), and PUTs the result back, so unspecified properties are preserved (a PATCH-like behavior).

Note that --routing-rules REPLACES the entire rule list; there is no way to edit a single rule in place. To modify one rule, fetch the current rules with 'ionosctl cdn ds rr get --distribution-id <id> -o json', edit the JSON, and pass the full array back. Provide 1-25 rules. See 'ionosctl cdn ds create --routing-rules-example' for the JSON format and field meanings (scheme, upstream host/caching/waf/rateLimitClass/sniMode/geoRestrictions).`,
		Example: `# Rebind the distribution to a new HTTPS certificate
ionosctl cdn ds update --distribution-id <id> --certificate-id 5a029f4a-72e5-11ec-90d6-0242ac120003

# Replace all routing rules from a file
ionosctl cdn ds update --distribution-id <id> --routing-rules rules.json`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagCDNDistributionID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			distributionId := viper.GetString(core.GetFlagName(c.NS, constants.FlagCDNDistributionID))
			r, _, err := client.Must().CDNClient.DistributionsApi.DistributionsFindById(context.Background(), distributionId).Execute()
			if err != nil {
				return fmt.Errorf("failed finding distribution: %w", err)
			}

			updated, err := updateDistribution(c, r)
			if err != nil {
				return err
			}

			return printDistribution(c, updated)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCDNDistributionID, constants.FlagIdShort, "", "The ID of the distribution you want to update",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.DistributionsProperty(func(r cdn.Distribution) string {
				return r.Id
			})
		}, constants.CDNApiRegionalURL, constants.CDNLocations),
	)
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return addDistributionCreateFlags(cmd)
}

func updateDistribution(c *core.CommandConfig, d cdn.Distribution) (cdn.Distribution, error) {
	input := &d.Properties
	err := setPropertiesFromFlags(c, input)
	if err != nil {
		return cdn.Distribution{}, err
	}

	rNew, _, err := client.Must().CDNClient.DistributionsApi.DistributionsPut(context.Background(), d.Id).
		DistributionUpdate(cdn.DistributionUpdate{Id: d.Id, Properties: *input}).Execute()
	if err != nil {
		return cdn.Distribution{}, err
	}

	return rNew, nil
}

func printDistribution(c *core.CommandConfig, d cdn.Distribution) error {
	return c.Printer(allCols).Print(d)
}
