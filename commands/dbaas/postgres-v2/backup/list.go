package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/viper"
)

func BackupListCmd() *core.Command {
	ctx := context.TODO()
	list := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "backup",
		Verb:       "list",
		Aliases:    []string{"ls"},
		ShortDesc:  "List PostgreSQL Backups",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Backups.",
		Example:    "ionosctl dbaas postgres-v2 backup list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunBackupList,
		InitClient: true,
	})
	list.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "Filter backups by Cluster ID",
		core.WithCompletion(func() []string {
			clusters, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersGet(context.Background()).Execute()
			if err != nil {
				return []string{}
			}
			return functional.Map(clusters.Items, func(c psqlv2.ClusterRead) string {
				return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
					c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)
			})
		}, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	list.AddInt32Flag(constants.FlagLimit, "", 100, "The limit of the number of items to return")
	list.AddInt32Flag(constants.FlagOffset, "", 0, "The offset of the listing")

	list.AddStringSliceFlag(constants.ArgCols, "", defaultBackupCols, table.ColsMessage(backupCols))

	return list
}

func RunBackupList(c *core.CommandConfig) error {
	req := client.Must().PostgresClientV2.BackupsApi.BackupsGet(context.Background())

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagClusterId)) {
		req = req.FilterClusterId(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLimit)) {
		req = req.Limit(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagLimit)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagOffset)) {
		req = req.Offset(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagOffset)))
	}

	backups, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))
	return c.Out(table.Sprint(backupCols, backups, cols, table.WithPrefix("items")))
}
