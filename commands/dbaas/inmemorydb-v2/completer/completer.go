package completer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
)

// Versions returns available In-Memory DB version strings for tab completion.
func Versions() []string {
	versions, _, err := client.Must().InMemoryDBClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v inmemorydb.VersionRead) string {
		if v.Properties.Version != nil {
			return *v.Properties.Version
		}
		return v.Id
	})
}

// VersionIds returns version IDs with version info for tab completion.
func VersionIds() []string {
	versions, _, err := client.Must().InMemoryDBClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v inmemorydb.VersionRead) string {
		ver := ""
		if v.Properties.Version != nil {
			ver = *v.Properties.Version
		}
		return fmt.Sprintf("%s\tv%s", v.Id, ver)
	})
}

// ClusterIds returns cluster IDs with descriptive info for tab completion.
func ClusterIds() []string {
	clusters, _, err := client.Must().InMemoryDBClientV2.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(clusters.Items, func(c inmemorydb.ClusterRead) string {
		return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
			c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)
	})
}

// SnapshotIds returns snapshot IDs with cluster info for tab completion.
func SnapshotIds() []string {
	snapshots, _, err := client.Must().InMemoryDBClientV2.SnapshotsApi.SnapshotsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(snapshots.Items, func(s inmemorydb.SnapshotRead) string {
		clusterId := ""
		if s.Properties.ClusterId != nil {
			clusterId = *s.Properties.ClusterId
		}
		return fmt.Sprintf("%s\tfor cluster '%s'", s.Id, clusterId)
	})
}

// SnapshotLocations returns snapshot location names for tab completion.
func SnapshotLocations() []string {
	locations, _, err := client.Must().InMemoryDBClientV2.SnapshotLocationsApi.SnapshotlocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l inmemorydb.SnapshotLocationRead) string {
		if l.Properties.Location != nil {
			return *l.Properties.Location
		}
		return l.Id
	})
}

// SnapshotLocationIds returns snapshot location IDs with location info for tab completion.
func SnapshotLocationIds() []string {
	locations, _, err := client.Must().InMemoryDBClientV2.SnapshotLocationsApi.SnapshotlocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l inmemorydb.SnapshotLocationRead) string {
		loc := ""
		if l.Properties.Location != nil {
			loc = *l.Properties.Location
		}
		return fmt.Sprintf("%s\t%s", l.Id, loc)
	})
}
