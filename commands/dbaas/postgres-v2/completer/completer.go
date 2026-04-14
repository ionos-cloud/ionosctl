package completer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
)

// PostgresVersions returns available PostgreSQL version strings for tab completion.
func PostgresVersions() []string {
	versions, _, err := client.Must().PostgresClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v psqlv2.PostgresVersionRead) string {
		if v.Properties.Version != nil {
			return *v.Properties.Version
		}
		return v.Id
	})
}

// VersionIds returns version IDs with version info for tab completion.
func VersionIds() []string {
	versions, _, err := client.Must().PostgresClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v psqlv2.PostgresVersionRead) string {
		ver := ""
		if v.Properties.Version != nil {
			ver = *v.Properties.Version
		}
		return fmt.Sprintf("%s\tv%s", v.Id, ver)
	})
}

// ClusterIds returns cluster IDs with descriptive info for tab completion.
func ClusterIds() []string {
	clusters, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(clusters.Items, func(c psqlv2.ClusterRead) string {
		return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
			c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)
	})
}

// BackupIds returns backup IDs with recovery time info for tab completion.
func BackupIds() []string {
	backups, _, err := client.Must().PostgresClientV2.BackupsApi.BackupsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	const timeFmt = "2006-01-02 15:04"
	return functional.Map(backups.Items, func(b psqlv2.BackupRead) string {
		latest := "now"
		if b.Properties.LatestRecoveryTargetTime != nil {
			latest = b.Properties.LatestRecoveryTargetTime.Time.Format(timeFmt)
		}
		earliest := "n/a"
		if b.Properties.EarliestRecoveryTargetTime != nil {
			earliest = b.Properties.EarliestRecoveryTargetTime.Time.Format(timeFmt)
		}
		return fmt.Sprintf("%s\tfor cluster '%s': earliest: '%s', latest: '%s'",
			b.Id, *b.Properties.ClusterId, earliest, latest)
	})
}

// BackupLocations returns backup location names for tab completion.
func BackupLocations() []string {
	locations, _, err := client.Must().PostgresClientV2.BackupLocationsApi.BackuplocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l psqlv2.BackupLocationRead) string {
		if l.Properties.Location != nil {
			return *l.Properties.Location
		}
		return l.Id
	})
}

// BackupLocationIds returns backup location IDs with location info for tab completion.
func BackupLocationIds() []string {
	locations, _, err := client.Must().PostgresClientV2.BackupLocationsApi.BackuplocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l psqlv2.BackupLocationRead) string {
		loc := ""
		if l.Properties.Location != nil {
			loc = *l.Properties.Location
		}
		return fmt.Sprintf("%s\t%s", l.Id, loc)
	})
}
