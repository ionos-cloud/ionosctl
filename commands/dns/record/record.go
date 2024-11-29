package record

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/spf13/viper"

	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	allCols     = []string{"Id", "Name", "Content", "Type", "Enabled", "FQDN", "ZoneId", "ZoneName", "State"}
	defaultCols = []string{"Id", "Name", "Content", "Type", "Enabled", "FQDN", "State"}
)

func RecordCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "record",
			Short:            "The sub-commands of 'ionosctl dns record' allow you to manage DNS records. Records allow directing traffic for a domain to its correct location.",
			Aliases:          []string{"r"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(RecordsGetCmd())
	cmd.AddCommand(ZonesRecordsDeleteCmd())
	cmd.AddCommand(ZonesRecordsPostCmd())
	cmd.AddCommand(ZonesRecordsFindByIdCmd())
	cmd.AddCommand(ZonesRecordsPutCmd())
	return cmd
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
	uid, errParseUuid := uuid.FromString(nameOrId)
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
			zoneId, err := utils.ZoneResolve(viper.GetString(fn))
			if err != nil {
				return req, err
			}
			req = req.FilterZoneId(zoneId)
		}

		if fn := core.GetFlagName(cmdNs, constants.FlagRecord); viper.IsSet(fn) {
			record := viper.GetString(fn)
			if _, ok := uuid.FromString(record); ok != nil /* not ok (name is provided) */ {
				req = req.FilterName(record)
			}
		}
		return req, nil
	}
}
