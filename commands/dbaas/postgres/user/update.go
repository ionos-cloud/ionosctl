package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UpdateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "update",
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Update a user's password",
			LongDesc: `Change the login password of an existing user in the given cluster.

Only the password can be changed; a user cannot be renamed. System users (System=true) are managed by the service and cannot be patched.

Required values to run command:

* Cluster Id
* User (name)
* Password (the new password)`,
			Example:   `ionosctl dbaas postgres user update --cluster-id CLUSTER_ID --user appuser --password 'N3wS3cr3tPass'`,
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster the user belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "Name of the existing (non-system) user whose password is changed")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgUser,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "New login password for the user. Minimum 10 characters")

	return c
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.ArgUser, constants.ArgPassword)
}

func runUpdateCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	username := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	password := viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))

	user, _, err := client.Must().PostgresClient.UsersApi.UsersPatch(
		context.Background(),
		clusterId,
		username,
	).UsersPatchRequest(psql.UsersPatchRequest{Properties: psql.PatchUserProperties{Password: &password}}).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(user)
}
