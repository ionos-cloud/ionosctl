package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

var secondaryZoneCompleterCols = []table.Column{
	{Name: "Id", JSONPath: "id"},
	{Name: "Name", JSONPath: "properties.zoneName"},
	{Name: "State", JSONPath: "metadata.state"},
}

func SecondaryZonesIDs() []string {
	secondaryZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	t := table.New(secondaryZoneCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(secondaryZones); err != nil {
		return nil
	}

	return completions.NewCompleter(
		t.Rows(), "Id",
	).AddInfo("Name").AddInfo("State", "(%v)").ToString()
}

// Zones returns all zones matching the given filters
func Zones(fs ...Filter) (ionoscloud.ZoneReadList, error) {
	req := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return ionoscloud.ZoneReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return ionoscloud.ZoneReadList{}, err
	}
	return ls, nil
}

func ZonesProperty[V any](f func(ionoscloud.ZoneRead) V, fs ...Filter) []V {
	recs, err := Zones(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

type Filter func(request ionoscloud.ApiZonesGetRequest) (ionoscloud.ApiZonesGetRequest, error)
