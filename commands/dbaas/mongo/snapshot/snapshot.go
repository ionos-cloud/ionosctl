package snapshot

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

func SnapshotCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Aliases:          []string{"snap", "backup", "snapshots", "backups"},
			Short:            "Mongo Snapshot Operations",
			Long:             "MongoDB Backups: A cluster can have multiple snapshots. A snapshot is added during the following cases:\nWhen a cluster is created, known as initial sync which usually happens in less than 24 hours.\nAfter a restore.\nA snapshot is a copy of the data in the cluster at a certain time. Every 24 hours, a base snapshot is taken, and every Sunday, a full snapshot is taken. Snapshots are retained for the last seven days; hence, recovery is possible for up to a week from the current date.\nYou can restore from any snapshot as long as it was created with the same or older MongoDB patch version.\nSnapshots are stored in an IONOS S3 Object Storage bucket in the same region as your database. Databases in regions where IONOS S3 Object Storage is not available is backed up to eu-central-2.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(SnapshotsListCmd())
	cmd.AddCommand(cluster.ClusterRestoreCmd())

	return cmd
}

var (
	allCols = []string{"SnapshotId", "CreationTime", "Size", "Version"}
)
