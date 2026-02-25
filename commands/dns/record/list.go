package record

import (
	"context"
	"fmt"
	"slices"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

var allColsSecondaryZoneRecord = append(
	slices.Clone(allCols[:len(allCols)-1]),
	table.Column{Name: "RootName", JSONPath: "metadata.rootName", Default: true},
)

func RecordsGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "dns",
			Resource:  "record",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "Retrieve all records from either a primary or secondary zone",
			Example: `ionosctl dns r list
ionosctl dns r list --secondary-zone SECONDARY_ZONE_ID
ionosctl dns r list --zone ZONE_ID`,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if c.Command.Command.Flags().Changed(constants.FlagZone) && c.Command.Command.Flags().Changed(constants.FlagSecondaryZone) {
					return fmt.Errorf("only one of the flags --%s and --%s can be set", constants.FlagZone, constants.FlagSecondaryZone)
				}

				if c.Command.Command.Flags().Changed(constants.FlagSecondaryZone) && c.Command.Command.Flags().Changed(constants.FlagName) {
					return fmt.Errorf("flag --%s is only available for zone records listing", constants.FlagName)
				}

				return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagZone}, []string{constants.FlagSecondaryZone}, []string{})
			},
			CmdRun: func(c *core.CommandConfig) error {
				if c.Command.Command.Flags().Changed(constants.FlagSecondaryZone) {
					return listSecondaryRecords(c)
				}

				return listRecordsCmd(c)
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", "(UUID or Zone Name) Filter used to fetch only the records that contain specified zone.")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ZonesProperty(
				func(t dns.ZoneRead) string {
					return t.Properties.ZoneName
				},
			), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Filter used to fetch only the records that contain specified record name. NOTE: Only available for zone records.")

	cmd.Command.Flags().String(constants.FlagSecondaryZone, "", "The name or ID of the secondary zone to fetch records from")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagSecondaryZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.SecondaryZonesIDs(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.PersistentFlags().StringSlice(
		constants.ArgCols, nil,
		fmt.Sprintf(
			"Set of columns to be printed on output \nAvailable columns for primary zones: %v\nAvailable columns for secondary zones: %v",
			table.AllCols(allCols), table.AllCols(allColsSecondaryZoneRecord),
		),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if cmd.Flags().Changed(constants.FlagSecondaryZone) {
				return table.AllCols(allColsSecondaryZoneRecord), cobra.ShellCompDirectiveNoFileComp
			}

			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func listRecordsCmd(c *core.CommandConfig) error {
	ls, err := Records(
		func(req dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
			if fn := core.GetFlagName(c.NS, constants.FlagZone); viper.IsSet(fn) {
				zoneId, err := utils.ZoneResolve(viper.GetString(fn))
				if err != nil {
					return req, err
				}
				req = req.FilterZoneId(zoneId)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				req = req.FilterName(viper.GetString(fn))
			}
			return req, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed listing zone records: %w", err)
	}

	items, ok := ls.GetItemsOk()
	if !ok || items == nil {
		return fmt.Errorf("could not retrieve Zone Record items")
	}

	t := table.New(allCols, table.WithPrefix("items"))
	if err := t.Extract(ls); err != nil {
		return err
	}

	for i, item := range items {
		if m, ok := item.GetMetadataOk(); ok && m != nil {
			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), m.ZoneId).Execute()
			if err == nil {
				t.SetCell(i, "ZoneName", z.Properties.ZoneName)
			}
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	out, err := t.Render(table.ResolveCols(allCols, cols))
	if err != nil {
		return err
	}
	fmt.Fprint(c.Command.Command.OutOrStdout(), out)
	return nil
}

func listSecondaryRecords(c *core.CommandConfig) error {
	records, err := secondaryRecords(c)
	if err != nil {
		return fmt.Errorf("failed listing secondary zone records: %w", err)
	}

	items, ok := records.GetItemsOk()
	if !ok || items == nil {
		return fmt.Errorf("could not retrieve Secondary Zone Record items")
	}

	t := table.New(allColsSecondaryZoneRecord, table.WithPrefix("items"))
	if err := t.Extract(records); err != nil {
		return err
	}

	for i, item := range items {
		if m, ok := item.GetMetadataOk(); ok && m != nil {
			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), m.ZoneId).Execute()
			if err == nil {
				t.SetCell(i, "ZoneName", z.Properties.ZoneName)
			}
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	out, err := t.Render(table.ResolveCols(allColsSecondaryZoneRecord, cols))
	if err != nil {
		return err
	}
	fmt.Fprint(c.Command.Command.OutOrStdout(), out)
	return nil
}

func secondaryRecords(c *core.CommandConfig) (dns.SecondaryZoneRecordReadList, error) {
	secondaryZoneIDOrName, _ := c.Command.Command.Flags().GetString(constants.FlagSecondaryZone)
	secondaryZoneID, err := utils.SecondaryZoneResolve(secondaryZoneIDOrName)
	if err != nil {
		return dns.SecondaryZoneRecordReadList{}, err
	}

	req := client.Must().DnsClient.RecordsApi.SecondaryzonesRecordsGet(context.Background(), secondaryZoneID)
	records, _, err := req.Execute()
	if err != nil {
		return dns.SecondaryZoneRecordReadList{}, err
	}

	return records, nil
}
