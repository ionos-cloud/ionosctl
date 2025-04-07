package utils

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
)

// SecondaryZoneResolve resolves nameOrId (the name of a zone, or the ID of a zone) - to the ID of the secondary zone.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
func SecondaryGatewayResolve(nameOrID string) (string, error) {
	if _, err := uuid.FromString(nameOrID); err == nil {
		return nameOrID, nil
	}

	secZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).FilterZoneName(nameOrID).Limit(1).Execute()
	//secZones1, _, err := client.Must().Apigateway.
	if err != nil {
		return "", fmt.Errorf("failed to retrieve zones by name %s: %w", nameOrID, err)
	}
	if secZones.Items == nil || len(secZones.Items) < 1 {
		return "", fmt.Errorf("no zones found with name %s", nameOrID)
	}

	return secZones.Items[0].Id, nil
}
