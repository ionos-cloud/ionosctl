package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/spf13/viper"
)

func SecondaryZonesIDs() []string {
	secondaryZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	secondaryZonesConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.DnsSecondaryZone, secondaryZones)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(
		secondaryZonesConverted, "Id",
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
