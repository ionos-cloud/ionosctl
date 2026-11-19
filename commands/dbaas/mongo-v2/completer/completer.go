package completer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
)

// Versions returns available MongoDB version strings for tab completion.
func Versions() []string {
	versions, _, err := client.Must().MongoClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v mongo.MongoVersionRead) string {
		if v.Properties.Version != nil {
			return *v.Properties.Version
		}
		return v.Id
	})
}

// VersionIds returns version IDs annotated with the version string for tab completion.
func VersionIds() []string {
	versions, _, err := client.Must().MongoClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(versions.Items, func(v mongo.MongoVersionRead) string {
		ver := ""
		if v.Properties.Version != nil {
			ver = *v.Properties.Version
		}
		return fmt.Sprintf("%s\tv%s", v.Id, ver)
	})
}

// ClusterIds returns cluster IDs with descriptive info for tab completion.
func ClusterIds() []string {
	clusters, _, err := client.Must().MongoClientV2.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(clusters.Items, func(c mongo.ClusterRead) string {
		switch p := c.Properties; {
		case p.EnterpriseCluster != nil:
			return fmt.Sprintf("%s\t%s: ENTERPRISE, %d instances, datacenter: %s",
				c.Id, p.EnterpriseCluster.Name, p.EnterpriseCluster.Instances.Count, p.EnterpriseCluster.Connection.DatacenterId)
		case p.BusinessCluster != nil:
			return fmt.Sprintf("%s\t%s: BUSINESS, %d instances, datacenter: %s",
				c.Id, p.BusinessCluster.Name, p.BusinessCluster.InstanceCount, p.BusinessCluster.Connection.DatacenterId)
		}
		return c.Id
	})
}

// SnapshotIds returns snapshot IDs with cluster info for tab completion.
func SnapshotIds() []string {
	snapshots, _, err := client.Must().MongoClientV2.SnapshotsApi.SnapshotsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(snapshots.Items, func(s mongo.SnapshotRead) string {
		clusterId := ""
		if s.Properties.ClusterId != nil {
			clusterId = *s.Properties.ClusterId
		}
		return fmt.Sprintf("%s\tfor cluster '%s'", s.Id, clusterId)
	})
}

// SnapshotLocations returns snapshot location names for tab completion.
func SnapshotLocations() []string {
	locations, _, err := client.Must().MongoClientV2.SnapshotLocationsApi.SnapshotlocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l mongo.SnapshotLocationRead) string {
		if l.Properties.Location != nil {
			return *l.Properties.Location
		}
		return l.Id
	})
}

// SnapshotLocationIds returns snapshot location IDs with location info for tab completion.
func SnapshotLocationIds() []string {
	locations, _, err := client.Must().MongoClientV2.SnapshotLocationsApi.SnapshotlocationsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(locations.Items, func(l mongo.SnapshotLocationRead) string {
		loc := ""
		if l.Properties.Location != nil {
			loc = *l.Properties.Location
		}
		return fmt.Sprintf("%s\t%s", l.Id, loc)
	})
}

// TemplateIds returns instance template IDs annotated with their sizing for tab
// completion. An empty edition returns the templates of every edition.
func TemplateIds(edition mongo.MongoEdition) []string {
	req := client.Must().MongoClientV2.InstanceTemplatesApi.InstancetemplatesGet(context.Background())
	if edition != "" {
		req = req.FilterEdition(edition)
	}
	templates, _, err := req.Execute()
	if err != nil {
		return nil
	}

	return functional.Map(templates.Items, func(t mongo.InstanceTemplateRead) string {
		p := t.Properties
		return fmt.Sprintf("%s\t%s (%s): %d cores, %dGB RAM, %dGB storage",
			t.Id, p.GetName(), p.GetEdition(), p.GetCores(), p.GetRam(), p.GetStorageSize())
	})
}
