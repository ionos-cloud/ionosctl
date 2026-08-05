package backup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func BackupListCmd() *core.Command {
	ctx := context.TODO()
	listEnv := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mariadb-v2",
		Resource:   "backup",
		Verb:       "list",
		Aliases:    []string{"ls", "l"},
		ShortDesc:  "List MariaDB Backups",
		LongDesc:   "Use this command to retrieve a list of MariaDB Backups. You can filter by cluster ID using `--cluster-id`.\n\nEach backup is a recovery WINDOW: the EarliestRecoveryTargetTime column shows the start of the range you can restore to (the window extends to the present). Zoom into a specific point with `--recovery-time` on `cluster restore` / `cluster create --backup-id`.",
		Example:    "ionosctl dbaas mariadb-v2 backup list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunBackupList,
		InitClient: true,
	})
	listEnv.AddStringFlag(constants.FlagClusterId, "", "", "Response filter to list only the backups of the specified cluster")
	listEnv.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	listEnv.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return listEnv
}

func RunBackupList(c *core.CommandConfig) error {
	return c.ListAllLocations(backupCols, func(cfg *shared.Configuration) (any, error) {
		apiClient := mariadb.NewAPIClient(cfg)
		req := apiClient.BackupsApi.BackupsGet(context.Background())

		if fn := core.GetFlagName(c.NS, constants.FlagClusterId); viper.IsSet(fn) {
			req = req.FilterClusterId(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
			req = req.Limit(viper.GetInt32(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
			req = req.Offset(viper.GetInt32(fn))
		}

		backups, _, err := req.Execute()
		return backups, err
	})
}
