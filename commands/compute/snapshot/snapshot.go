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
	{Name: "LicenceType", JSONPath: "properties.licenseType", Default: true},
	{Name: "Size", JSONPath: "properties.size", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func SnapshotCmd() *core.Command {
	snapshotCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Aliases:          []string{"ss", "snap"},
			Short:            "Snapshot Operations",
			Long:             "The sub-commands of `ionosctl compute snapshot` allow you to see information, to create, update, delete Snapshots.",
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
