package version

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func VersionListCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "version",
		Verb:       "list",
		Aliases:    []string{"ls"},
		ShortDesc:  "List PostgreSQL Versions",
		LongDesc:   "Use this command to retrieve a list of available PostgreSQL Versions.",
		Example:    "ionosctl dbaas postgres version list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunVersionList,
		InitClient: true,
	})
	cmd.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return cmd
}

func RunVersionList(c *core.CommandConfig) error {
	req := client.Must().PostgresClientV2.VersionsApi.VersionsGet(context.Background())
	if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
		req = req.Limit(viper.GetInt32(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
		req = req.Offset(viper.GetInt32(fn))
	}
	versions, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(versionCols, versions, cols, table.WithPrefix("items")))
}
