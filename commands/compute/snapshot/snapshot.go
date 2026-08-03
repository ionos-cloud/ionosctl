package snapshot

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allSnapshotCols = []table.Column{
	{Name: "SnapshotId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "LicenceType", JSONPath: "properties.licenceType", Default: true},
	{Name: "Size", JSONPath: "properties.size", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func SnapshotCmd() *core.Command {
	snapshotCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "snapshot",
			Aliases: []string{"ss", "snap"},
			Short:   "Manage Snapshots (point-in-time images of Compute Engine storage Volumes)",
			Long: `A Snapshot is a point-in-time image of a single Block Storage Volume, stored independently of the Volume it was taken from. Once created it lives on its own (it is not deleted when the source Volume is deleted) and can be used two ways: to restore a Volume back to the captured state (` + "`snapshot restore`" + `), or as a boot template when creating a new Volume (pass the Snapshot Id as the image when creating a Volume).

Scope: a Snapshot is tied to the physical LOCATION (region) where its source Volume lives, and to your contract. It can be used across multiple Virtual Data Centers, but only ones in that same location, and only within your own contract. To move data to another region, replicate it with the IONOS Backup Service or Object Storage instead - that is the practical difference between a Snapshot (fast, same-location clone/rollback) and a backup unit (` + "`ionosctl backupunit`" + `, cross-location redundancy for the IONOS Backup Service agent).

Size: a Snapshot captures the FULL provisioned capacity of the source Volume, including empty space. A 100 GB Volume with 10 GB of data still produces a 100 GB Snapshot.

Each Snapshot records the source OS licence type and a set of hot-plug capabilities; a Volume created or restored from it starts out with those same properties (see ` + "`snapshot update`" + ` to adjust them).`,
			TraverseChildren: true,
		},
	}
	snapshotCmd.AddColsFlag(allSnapshotCols)

	snapshotCmd.AddCommand(SnapshotListCmd())
	snapshotCmd.AddCommand(SnapshotGetCmd())
	snapshotCmd.AddCommand(SnapshotCreateCmd())
	snapshotCmd.AddCommand(SnapshotUpdateCmd())
	snapshotCmd.AddCommand(SnapshotRestoreCmd())
	snapshotCmd.AddCommand(SnapshotDeleteCmd())

	return core.WithConfigOverride(snapshotCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
