package reverse_record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "reverse-record",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a record. Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-and-manage-reverse-dns",
		Example:   "ionosctl dns rr create --name mail.example.com --ips 5.6.7.8",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagIps); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			rec, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsPut(context.Background(), uuidgen.Must()).
				ReverseRecordEnsure(dns.ReverseRecordEnsure{
					Properties: &dns.ReverseRecord{
						Name:        pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagName))),
						Ip:          pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagIp))),
						Description: pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))),
					},
				}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsRecord, rec,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the DNS Reverse Record.", core.RequiredFlagOption())
	cmd.AddStringSliceFlag(constants.FlagIps, "", []string{}, "[IPv4/IPv6] Specifies for which IP address the reverse record should be created. The IP addresses needs to be owned by the contract", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description stored along with the reverse DNS record to describe its usage")
	cmd.Command.SilenceUsage = true

	return cmd
}
