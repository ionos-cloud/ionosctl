package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	dns "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
)

func RecordsGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "Retrieve all records",
		Example:   "ionosctl dns record list",
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := Records(func(req dns.ApiRecordsGetRequest) dns.ApiRecordsGetRequest {
				if fn := core.GetFlagName(c.NS, constants.FlagZoneId); viper.IsSet(fn) {
					req = req.FilterZoneId(viper.GetString(fn))
				}
				if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
					req = req.FilterName(viper.GetString(fn))
				}
				if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
					req = req.Offset(viper.GetInt32(fn))
				}
				if fn := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(fn) {
					req = req.Limit(viper.GetInt32(fn))
				}
				return req
			})

			if err != nil {
				return fmt.Errorf("failed listing records: %w", err)
			}

			return c.Printer.Print(getRecordsPrint(c, *ls))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZoneId, "", "", "Filter used to fetch only the records that contain specified zoneId")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZoneIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagName, "", "", "Filter used to fetch only the records that contain specified record name")
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "The first element (of the total list of elements) to include in the response. Use together with limit for pagination")
	cmd.AddInt32Flag(constants.FlagMaxResults, "", 0, constants.DescMaxResults)

	return cmd
}
