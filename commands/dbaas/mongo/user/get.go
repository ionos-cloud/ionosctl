package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a MongoDB user",
		Example:   "ionosctl dbaas mongo user get",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(FlagDatabase)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.ArgUser)
			if err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			user := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting user %s...", user))

			u, _, err := client.Must().MongoClient.UsersApi.
				ClustersUsersFindById(context.Background(), clusterId, user).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			uConverted, err := resource2table.ConvertDbaasMongoUserToTable(u)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(u, uConverted, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(FlagDatabase, FlagDatabaseShort, "", "The authentication database")
	cmd.AddStringFlag(constants.ArgUser, "", "", "The authentication username")

	cmd.Command.SilenceUsage = true

	return cmd
}
