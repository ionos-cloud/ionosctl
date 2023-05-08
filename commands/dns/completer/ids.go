package completer

import (
	"context"

	dns "github.com/ionos-cloud/sdk-go-dnsaas"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
)

func Zones() []string {
	ls, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t dns.ZoneResponse) string {
		return *t.GetId()
	})
}

type RecordFilter func(dns.ApiRecordsGetRequest) dns.ApiRecordsGetRequest

func Records(filters ...RecordFilter) []string {
	req := client.Must().DnsClient.RecordsApi.RecordsGet(context.Background())

	for _, f := range filters {
		req = f(req)
	}

	ls, _, err := req.Execute()
	if err != nil {
		return nil
	}

	return functional.Map(*ls.GetItems(), func(t dns.RecordResponse) string {
		return *t.GetId()
	})
}
