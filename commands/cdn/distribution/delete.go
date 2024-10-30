package distribution

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cdn "github.com/ionos-cloud/sdk-go-cdn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "distribution",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a distribution",
		Example:   `ionosctl cdn ds delete --distribution-id ID`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagCDNDistributionID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			distributionID := viper.GetString(core.GetFlagName(c.NS, constants.FlagCDNDistributionID))
			d, _, err := client.Must().CDNClient.DistributionsApi.DistributionsFindById(context.Background(), distributionID).Execute()
			if err != nil {
				return fmt.Errorf("distribution not found: %w", err)
			}

			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete distribution %s for domain %s", *d.Id, *d.Properties.Domain),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf("user cancelled deletion")
			}

			_, err = client.Must().CDNClient.DistributionsApi.DistributionsDelete(context.Background(), *d.Id).Execute()
			return err
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
