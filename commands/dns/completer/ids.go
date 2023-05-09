package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/record"

	dns "github.com/ionos-cloud/sdk-go-dnsaas"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
)

func ZoneIds() []string {
	ls, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t dns.ZoneResponse) string {
		return *t.GetId()
	})
}

func RecordIds(filters ...record.Filter) []string {
	ls, err := record.Records(filters...)
	if err != nil {
		return nil
	}

	return functional.Map(*ls.GetItems(), func(t dns.RecordResponse) string {
		return *t.GetId()
	})
}
