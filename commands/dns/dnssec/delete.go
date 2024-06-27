package dnssec

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "dnssec",
		Verb:      "delete",
		Aliases:   []string{"del", "rm", "remove"},
		ShortDesc: "Removes ALL associated DNSKEY records for your DNS zone and disables DNSSEC keys.",
		Example:   `ionosctl dns keys delete --zone ZONE`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			zoneName := viper.GetString(core.GetFlagName(c.NS, constants.FlagZone))
			zoneId, err := zone.Resolve(zoneName)
			if err != nil {
				return err
			}

			_, _, err = client.Must().DnsClient.DNSSECApi.ZonesKeysDelete(context.Background(), zoneId).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "DNSKEY records deleted and DNSSEC keys disabled for zone %s\n", zoneName)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZonesProperty(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
