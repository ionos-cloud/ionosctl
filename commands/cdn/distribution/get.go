package distribution

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	cdn "github.com/ionos-cloud/sdk-go-cdn"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func FindByID() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "distribution",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a distribution",
		Example:   "ionosctl cdn ds get --distribution-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagCDNDistributionID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			distributionID := viper.GetString(core.GetFlagName(c.NS, constants.FlagCDNDistributionID))
			r, _, err := client.Must().CDNClient.DistributionsApi.DistributionsFindById(context.Background(),
				distributionID,
			).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.CDNDistribution, r,
				tabheaders.GetHeadersAllDefault(defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagCDNDistributionID, constants.FlagIdShort, "", "The ID of the distribution you want to retrieve", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCDNDistributionID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DistributionsProperty(func(r cdn.Distribution) string {
			return *r.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
