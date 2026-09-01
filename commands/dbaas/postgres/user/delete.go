package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
			LongDesc:  `Delete a user (login role) from the given cluster. System users cannot be deleted. Databases owned by the user are not dropped; reassign or drop them separately if needed.` + "\n\n" + `Required values to run command:` + "\n\n" + `* Cluster Id` + "\n" + `* User (name)`,
			Example:   `ionosctl dbaas postgres user delete --cluster-id CLUSTER_ID --user appuser`,
			PreCmdRun: preRunDeleteCmd,
			CmdRun:    runDeleteCmd,
		},
	)
	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster the user belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "Name of the (non-system) user to delete")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgUser,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveNoFileComp
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

	c.Msg("DbaaS Postgres User %v successfully deleted", username)
	return nil
}
