package snapshot

import (
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
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

type SnapshotPrint struct {
	SnapshotId   string    `json:"SnapshotId,omitempty"`
	CreationTime time.Time `json:"CreationTime,omitempty"`
	Size         int32     `json:"Size,omitempty"`
	Version      string    `json:"Version,omitempty"`
}

var (
	allJSONPaths = map[string]string{
		"SnapshotId":   "id",
		"CreationTime": "properties.creationTime",
		"Size":         "properties.size",
		"Version":      "properties.version",
	}

	allCols = []string{"SnapshotId", "CreationTime", "Size", "Version"}
)

func MakeSnapshotPrintObject(snapshots *[]ionoscloud.SnapshotResponse) []map[string]interface{} {
	if snapshots == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(*snapshots))
	for _, snapshot := range *snapshots {
		var snapshotPrint SnapshotPrint

		if snapshot.GetId() != nil {
			snapshotPrint.SnapshotId = *snapshot.GetId()
		}

		properties := snapshot.GetProperties()
		if properties != nil {
			if properties.GetCreationTime() != nil {
				snapshotPrint.CreationTime = *properties.GetCreationTime()
			}
			if properties.GetSize() != nil {
				snapshotPrint.Size = *properties.GetSize()
			}
			if properties.GetVersion() != nil {
				snapshotPrint.Version = *properties.GetVersion()
			}
		}

		o := structs.Map(snapshotPrint)
		out = append(out, o)
	}

	return out
}

func getSnapshotPrint(c *core.CommandConfig, dcs *[]ionoscloud.SnapshotResponse) printer.Result {
	r := printer.Result{}
	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = MakeSnapshotPrintObject(dcs)                                                                                                 // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(structs.Names(SnapshotPrint{}), viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}
