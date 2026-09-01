package reverse_record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "reverse-record",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a reverse DNS (PTR) record",
		LongDesc: `Create a reverse DNS record so a PTR lookup on --ip returns --name.

--ip must be an IPv4 or IPv6 address owned by your contract (e.g. from a reserved IP block); --name is the hostname it should resolve back to. Commonly used to give a mail server matching forward and reverse DNS.

Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-and-manage-reverse-dns`,
		Example: "ionosctl dns reverse-record create --name mail.example.com --ip 5.6.7.8",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagIp); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			rec, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsPut(context.Background(), uuidgen.Must()).
				ReverseRecordEnsure(dns.ReverseRecordEnsure{
					Properties: dns.ReverseRecord{
						Name:        viper.GetString(core.GetFlagName(c.NS, constants.FlagName)),
						Ip:          viper.GetString(core.GetFlagName(c.NS, constants.FlagIp)),
						Description: pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))),
					},
				}).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(rec)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Hostname the IP should resolve back to, e.g. mail.example.com", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagIp, "", "", "IPv4 or IPv6 address to create the reverse record for; must be owned by your contract", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagIp, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ipblocks, _, err := client.Must().CloudClient.IPBlocksApi.IpblocksGet(context.Background()).Execute()
		if err != nil || ipblocks.Items == nil || len(*ipblocks.Items) == 0 {
			return nil, cobra.ShellCompDirectiveError
		}
		var ips []string
		for _, ipblock := range *ipblocks.Items {
			if ipblock.Properties.Ips != nil {
				ips = append(ips, *ipblock.Properties.Ips...)
			}
		}
		return ips, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagDescription, "", "", "Free-text note for your own reference; not served in DNS")
	cmd.Command.SilenceUsage = true

	return cmd
}
