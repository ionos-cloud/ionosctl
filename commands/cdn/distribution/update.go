package distribution

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "distribution",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a distribution's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl cdn ds update --distribution-id",
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
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	out, err := jsontabwriter.GenerateOutput("", jsonpaths.CDNDistribution, d,
		tabheaders.GetHeadersAllDefault(defaultCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
