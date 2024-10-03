package utils

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-dns"

	"github.com/spf13/viper"
)

// SecondaryZoneResolve resolves nameOrId (the name of a zone, or the ID of a zone) - to the ID of the secondary zone.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
func SecondaryZoneResolve(nameOrID string) (string, error) {
	if _, err := uuid.FromString(nameOrID); err == nil {
		return nameOrID, nil
	}

	secZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).FilterZoneName(nameOrID).Limit(1).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve zones by name %s: %w", nameOrID, err)
	}
	if secZones.Items == nil || len(*secZones.Items) < 1 {
		return "", fmt.Errorf("no zones found with name %s", nameOrID)
	}

	return *(*secZones.Items)[0].Id, nil
}

// ZoneResolve resolves nameOrId (the name of a zone, or the ID of a zone) - to the ID of the zone.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
func ZoneResolve(nameOrId string) (string, error) {
	uid, errParseUuid := uuid.FromString(nameOrId)
	zId := uid.String()
	if errParseUuid != nil {
		// nameOrId is a name
		ls, _, errFindZoneByName := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).FilterZoneName(nameOrId).Limit(1).Execute()
		if errFindZoneByName != nil {
			return "", fmt.Errorf("failed finding a zone by name: %w", errFindZoneByName)
		}
		if len(*ls.Items) < 1 {
			return "", fmt.Errorf("could not find zone by name %s: got %d zones", nameOrId, len(*ls.Items))
		}
		zId = *(*ls.Items)[0].Id
	}
	return zId, nil
}

// Zones returns all zones matching the given filters
func Zones(fs ...Filter) (ionoscloud.ZoneReadList, error) {
	// Hack to enforce the dns-level flag default for API URL on the completions too
	if url := config.GetServerUrl(); url == constants.DefaultApiURL {
		viper.Set(constants.ArgServerUrl, "")
	}

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
	return functional.Map(*recs.Items, f)
}

type Filter func(request ionoscloud.ApiZonesGetRequest) (ionoscloud.ApiZonesGetRequest, error)
