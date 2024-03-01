package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "delete",
			Aliases:   []string{"del"},
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Delete user",
			LongDesc:  "Delete the specified user from the given database cluster",
			Example:   "ionosctl dbaas-postgres user delete --cluster-id <cluster-id> --user <user>",
			PreCmdRun: preRunDeleteCmd,
			CmdRun:    runDeleteCmd,
		},
	)
	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveDefault
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "The name of the user to delete")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgUser,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveDefault
		},
	)

	return c
}

func preRunDeleteCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.ArgUser)
}

func runDeleteCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	username := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))

	if !confirm.FAsk(
		c.Command.Command.InOrStdin(), fmt.Sprintf("delete user %s from cluster %s", username, clusterId),
		viper.GetBool(constants.ArgForce),
	) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := client.Must().PostgresClient.UsersApi.UsersDelete(context.Background(), clusterId, username).Execute()
	if err != nil {
		return err
	}

	out := jsontabwriter.GenerateLogOutput("DbaaS Postgres User successfully deleted")

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
