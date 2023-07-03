package zone

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a zone",
		Example:   "ionosctl dns z delete --zone ZONE",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagZone})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			zoneId, err := Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), zoneId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting zone by id %s: %w", zoneId, err)
			}
			yes := confirm.Ask(fmt.Sprintf("Are you sure you want to delete zone %s (%s)", *z.Properties.ZoneName, *z.Properties.Description),
				viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
			if !yes {
				return nil
			}

			_, err = client.Must().DnsClient.ZonesApi.ZonesDelete(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)),
			).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", fmt.Sprintf("%s. Required or -%s", constants.DescZone, constants.ArgAllShort))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return ZonesProperty(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no confirmation")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all zones. Required or -%s", constants.FlagZoneShort))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Deleting all zones!")
	xs, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(c.Context).Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(*xs.GetItems(), func(z dns.ZoneRead) error {
		yes := confirm.Ask(fmt.Sprintf("Are you sure you want to delete zone %s (desc: %s)", *z.Properties.ZoneName, *z.Properties.Description),
			viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
		if yes {
			_, delErr := client.Must().DnsClient.ZonesApi.ZonesDelete(c.Context, *z.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", *z.Id, *z.Properties.ZoneName, delErr)
			}
		}
		return nil
	})

	return err
}
