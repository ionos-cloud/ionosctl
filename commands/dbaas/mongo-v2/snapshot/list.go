package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func SnapshotListCmd() *core.Command {
	ctx := context.TODO()
	listEnv := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mongo-v2",
		Resource:   "snapshot",
		Verb:       "list",
		Aliases:    []string{"ls", "l"},
		ShortDesc:  "List MongoDB Snapshots",
		LongDesc:   "Use this command to retrieve a list of MongoDB Snapshots. You can filter by cluster ID using `--cluster-id`.\n\nEach snapshot is a recovery WINDOW: the EarliestRecoveryTargetTime and LatestRecoveryTargetTime columns bound the range you can restore to (LatestRecoveryTargetTime is empty while the window extends to the present). Pick a specific point in the window with `--recovery-time` on `cluster restore` / `cluster create --snapshot-id`.",
		Example:    "ionosctl dbaas mongo-v2 snapshot list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunSnapshotList,
		InitClient: true,
	})
	listEnv.AddStringFlag(constants.FlagClusterId, "", "", "Response filter to list only the snapshots of the specified cluster")
	listEnv.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	listEnv.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return listEnv
}

func RunSnapshotList(c *core.CommandConfig) error {
	return c.ListAllLocations(snapshotCols, func(cfg *shared.Configuration) (any, error) {
		apiClient := mongo.NewAPIClient(cfg)
		req := apiClient.SnapshotsApi.SnapshotsGet(context.Background())

		if fn := core.GetFlagName(c.NS, constants.FlagClusterId); viper.IsSet(fn) {
			req = req.FilterClusterId(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
			req = req.Limit(viper.GetInt32(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
			req = req.Offset(viper.GetInt32(fn))
		}

		snapshots, _, err := req.Execute()
		return snapshots, err
	})
}
