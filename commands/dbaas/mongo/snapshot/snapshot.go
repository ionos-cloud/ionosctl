package snapshot

import (
	"time"

	"github.com/fatih/structs"
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
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(SnapshotsListCmd())

	return cmd
}

type SnapshotPrint struct {
	SnapshotId   string    `json:"SnapshotId,omitempty"`
	CreationTime time.Time `json:"CreationTime,omitempty"`
	Size         int32     `json:"Size,omitempty"`
	Version      string    `json:"Version,omitempty"`
}

var allCols = structs.Names(SnapshotPrint{})

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
