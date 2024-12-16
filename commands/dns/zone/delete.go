package zone

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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

			zoneId, err := utils.ZoneResolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), zoneId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting zone by id %s: %w", zoneId, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete zone %s (desc: '%s')", *z.Properties.ZoneName, *z.Properties.Description),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, _, err = client.Must().DnsClient.ZonesApi.ZonesDelete(context.Background(), zoneId).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", fmt.Sprintf("%s. Required or -%s", constants.DescZone, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			return completer.ZonesProperty(func(t dns.ZoneRead) string {
				return *t.Properties.ZoneName
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all zones. Required or -%s", constants.FlagZoneShort))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all zones!"))
	xs, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(c.Context).Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(*xs.GetItems(), func(z dns.ZoneRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete zone %s (desc: '%s')", *z.Properties.ZoneName, *z.Properties.Description),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, _, delErr := client.Must().DnsClient.ZonesApi.ZonesDelete(c.Context, *z.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", *z.Id, *z.Properties.ZoneName, delErr)
			}
		}
		return nil
	})

	return err
}
