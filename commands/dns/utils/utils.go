package utils

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
