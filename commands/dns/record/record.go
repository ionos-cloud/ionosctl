package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/spf13/viper"

	"github.com/google/uuid"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
)

func RecordCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "record",
			Short:            "DNS RecordsProperty",
			Aliases:          []string{"r"},
			Long:             "The sub-commands of `ionosctl dns record` allow you to perform operations on DNS records",
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	cmd.AddCommand(RecordsGetCmd())
	cmd.AddCommand(ZonesRecordsDeleteCmd())
	cmd.AddCommand(ZonesRecordsPostCmd())
	cmd.AddCommand(ZonesRecordsFindByIdCmd())
	cmd.AddCommand(ZonesRecordsPutCmd())
	return cmd
}

// Helper functions for printing record

func getRecordsPrint(c *core.CommandConfig, data dns.RecordReadList) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		r.OutputJSON = data.Items // TODO: See above comment. Remove `.Items` once JSON marshalling works as one would expect
		r.KeyValue = makeRecordPrintObj(*data.Items...)
		r.Columns = printer.GetHeaders(allCols, defaultCols, cols)
	}
	return r
}

func getRecordPrint(c *core.CommandConfig, data dns.RecordRead) printer.Result {
	return getRecordsPrint(c, dns.RecordReadList{Items: &[]dns.RecordRead{data}})
}

type recordPrint struct {
	Id       string `json:"ID,omitempty"`
	Name     string `json:"Name,omitempty"`
	Content  string `json:"Content,omitempty"`
	Type     string `json:"Type,omitempty"`
	Enabled  bool   `json:"Enabled,omitempty"`
	FQDN     string `json:"FQDN,omitempty"`
	State    string `json:"State,omitempty"`
	ZoneId   string `json:"ZoneId,omitempty"`
	ZoneName string `json:"ZoneName,omitempty"`
}

var allCols = structs.Names(recordPrint{})
var defaultCols = allCols[:len(allCols)-2]

func makeRecordPrintObj(data ...dns.RecordRead) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(data))

	for _, item := range data {
		var printObj recordPrint
		printObj.Id = *item.GetId()

		// Fill in the rest of the fields from the response object

		if propertiesOk, ok := item.GetPropertiesOk(); ok && propertiesOk != nil {
			printObj.Type = string(*propertiesOk.Type)

			printObj.Enabled = *propertiesOk.Enabled
			printObj.Content = *propertiesOk.Content
			printObj.Name = *propertiesOk.Name
		}
		if m, ok := item.GetMetadataOk(); ok && m != nil {
			printObj.FQDN = *m.Fqdn
			printObj.State = string(*m.State)
			printObj.ZoneId = *m.ZoneId
			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), *m.ZoneId).Execute()
			if err == nil && z.Properties != nil {
				printObj.ZoneName = *z.Properties.ZoneName
			}
		}

		o := structs.Map(printObj)
		out = append(out, o)
	}
	return out
}

// RecordsProperty returns a list of properties of all records matching the given filters
func RecordsProperty[V any](f func(dns.RecordRead) V, fs ...Filter) []V {
	recs, err := Records(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

// Records returns all records matching the given filters
func Records(fs ...Filter) (dns.RecordReadList, error) {
	req := client.Must().DnsClient.RecordsApi.RecordsGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return dns.RecordReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return dns.RecordReadList{}, err
	}
	return ls, nil
}

// Resolve resolves nameOrId (the name of a record, or the ID of a record) - to the ID of the record.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
func Resolve(nameOrId string) (string, error) {
	uid, errParseUuid := uuid.Parse(nameOrId)
	rId := uid.String()
	if errParseUuid != nil {
		// nameOrId is a name
		ls, _, err := client.Must().DnsClient.RecordsApi.RecordsGet(context.Background()).FilterName(nameOrId).Limit(1).Execute()
		if err != nil {
			return "", fmt.Errorf("failed finding a record by name %s: %w", nameOrId, err)
		}
		if len(*ls.Items) < 1 {
			return "", fmt.Errorf("could not find record by name %s: got %d records", nameOrId, len(*ls.Items))
		}
		rId = *(*ls.Items)[0].Id
	}
	return rId, nil
}

type Filter func(dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error)

func FilterRecordsByZoneAndRecordFlags(cmdNs string) Filter {
	return func(req dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
		if fn := core.GetFlagName(cmdNs, constants.FlagZone); viper.IsSet(fn) {
			zoneId, err := zone.Resolve(viper.GetString(fn))
			if err != nil {
				return req, err
			}
			req = req.FilterZoneId(zoneId)
		}

		if fn := core.GetFlagName(cmdNs, constants.FlagRecord); viper.IsSet(fn) {
			record := viper.GetString(fn)
			if _, ok := uuid.Parse(record); ok != nil /* not ok (name is provided) */ {
				req = req.FilterName(record)
			}
		}
		return req, nil
	}
}
