package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
)

func List() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "backup",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List MariaDB Backups",
		LongDesc:  "Use this command to retrieve a list of MariaDB Backups provisioned under your account.",
		Example:   "ionosctl dbaas mariadb backup list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Clusters..."))

			backups, err := Backups(FilterPaginationFlags(c))
			if err != nil {
				return err
			}

			// cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			// out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DbaasPostgresBackup,
			// 	backups, tabheaders.GetHeaders(allCols, defaultCols, cols))
			// if err != nil {
			// 	return err
			// }
			j, err := backups.MarshalJSON()
			fmt.Fprintf(c.Command.Command.OutOrStdout(), string(j))
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the MariaDB Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	return cmd
}
