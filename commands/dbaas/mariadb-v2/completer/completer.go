package completer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
)

// Versions returns available MariaDB version strings for tab completion.
func Versions() []string {
	versions, _, err := client.Must().MariaClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v mariadb.MariadbVersionRead) string {
		if v.Properties.Version != nil {
			return *v.Properties.Version
		}
		return v.Id
	})
}

// VersionIds returns version IDs annotated with the version string for tab completion.
func VersionIds() []string {
	versions, _, err := client.Must().MariaClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v mariadb.MariadbVersionRead) string {
		ver := ""
		if v.Properties.Version != nil {
			ver = *v.Properties.Version
		}
		return fmt.Sprintf("%s\tv%s", v.Id, ver)
	})
}

// ClusterIds returns cluster IDs with descriptive info for tab completion.
func ClusterIds() []string {
	clusters, _, err := client.Must().MariaClientV2.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(clusters.Items, func(c mariadb.ClusterRead) string {
		return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
			c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)
	})
}

// BackupIds returns backup IDs with cluster info for tab completion.
func BackupIds() []string {
	backups, _, err := client.Must().MariaClientV2.BackupsApi.BackupsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(backups.Items, func(b mariadb.BackupRead) string {
		clusterId := ""
		if b.Properties.ClusterId != nil {
			clusterId = *b.Properties.ClusterId
		}
		return fmt.Sprintf("%s\tfor cluster '%s'", b.Id, clusterId)
	})
}

// BackupLocations returns backup location names for tab completion.
func BackupLocations() []string {
	locations, _, err := client.Must().MariaClientV2.BackupLocationsApi.BackuplocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l mariadb.BackupLocationRead) string {
		if l.Properties.Location != nil {
			return *l.Properties.Location
		}
		return l.Id
	})
}

// BackupLocationIds returns backup location IDs with location info for tab completion.
func BackupLocationIds() []string {
	locations, _, err := client.Must().MariaClientV2.BackupLocationsApi.BackuplocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l mariadb.BackupLocationRead) string {
		loc := ""
		if l.Properties.Location != nil {
			loc = *l.Properties.Location
		}
		return fmt.Sprintf("%s\t%s", l.Id, loc)
	})
}
