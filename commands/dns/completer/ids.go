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
