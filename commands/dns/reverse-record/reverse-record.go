package reverse_record

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/spf13/viper"

	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "IP", "Description"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "reverse-record",
			Short:            "The sub-commands of 'ionosctl dns reverse-record' allow you to manage DNS reverse records.",
			Aliases:          []string{"rr"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, allCols, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(List())
	cmd.AddCommand(Create())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Update())
	cmd.AddCommand(Get())
	return cmd
}

// RecordsProperty returns a list of properties of all records matching the given filters
func RecordsProperty[V any](f func(dns.ReverseRecordRead) V, fs ...Filter) []V {
	recs, err := Records(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

// Records returns all records matching the given filters
func Records(fs ...Filter) (dns.ReverseRecordsReadList, error) {
	req := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return dns.ReverseRecordsReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return dns.ReverseRecordsReadList{}, err
	}
	return ls, nil
}

// Resolve resolves ipOrId (the IP of a record, or the ID of a record) - to the ID of the record.
// If it's an ID, it's returned as is. If it's not, then it's an IP, and we try to resolve it
func Resolve(ipOrId string) (string, error) {
	uid, errParseUuid := uuid.FromString(ipOrId)
	rId := uid.String()
	if errParseUuid != nil {
		// nameOrId is a name
		ls, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsGet(context.Background()).FilterRecordIp([]string{ipOrId}).Limit(1).Execute()
		if err != nil {
			return "", fmt.Errorf("failed finding a record by IP %s: %w", ipOrId, err)
		}
		if len(ls.Items) < 1 {
			return "", fmt.Errorf("could not find record by IP %s: got %d records", ipOrId, len(ls.Items))
		}
		rId = ls.Items[0].Id
	}
	return rId, nil
}

type Filter func(request dns.ApiReverserecordsGetRequest) (dns.ApiReverserecordsGetRequest, error)

func FilterRecordsByIp(cmdNs string) Filter {
	return func(req dns.ApiReverserecordsGetRequest) (dns.ApiReverserecordsGetRequest, error) {
		if fn := core.GetFlagName(cmdNs, constants.FlagIps); viper.IsSet(fn) {
			req = req.FilterRecordIp(viper.GetStringSlice(fn))
		}

		return req, nil
	}
}

func FilterLimitOffset(cmdNs string) Filter {
	return func(req dns.ApiReverserecordsGetRequest) (dns.ApiReverserecordsGetRequest, error) {
		if fn := core.GetFlagName(cmdNs, constants.FlagOffset); viper.IsSet(fn) {
			req = req.Offset(viper.GetInt32(fn))
		}
		if fn := core.GetFlagName(cmdNs, constants.FlagMaxResults); viper.IsSet(fn) {
			req = req.Limit(viper.GetInt32(fn))
		}

		return req, nil
	}
}
