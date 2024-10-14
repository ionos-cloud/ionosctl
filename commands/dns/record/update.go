package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ZonesRecordsPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl dns z update --zone ZONE_ID --record RECORD",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone, constants.FlagRecord); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			zoneId, err := utils.ZoneResolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}
			recordId, err := Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecord)))
			if err != nil {
				return err
			}
			r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(), zoneId, recordId).Execute()
			if err != nil {
				return fmt.Errorf("failed finding record: %w", err)
			}
			return partiallyUpdateRecordAndPrint(c, r)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ZonesProperty(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagRecord, constants.FlagRecordShort, "", "The ID or name of the DNS record", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecord, func(cobraCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordsProperty(func(r dns.RecordRead) string {
			return *r.Properties.Name
		}, FilterRecordsByZoneAndRecordFlags(cmd.NS)), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return addRecordCreateFlags(cmd)
}

func partiallyUpdateRecordAndPrint(c *core.CommandConfig, r dns.RecordRead) error {
	input := r.Properties
	modifyRecordPropertiesFromFlags(c, input)

	rNew, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsPut(context.Background(), *r.Metadata.ZoneId, *r.Id).
		RecordEnsure(dns.RecordEnsure{Properties: input}).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	// if err != nil {
	//	return err
	// }

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsRecord, rNew,
		tabheaders.GetHeadersAllDefault(defaultCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
