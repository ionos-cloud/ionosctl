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

func CreateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Create a user",
			LongDesc: `Create a new PostgreSQL login role in the given cluster.

The cluster must already be AVAILABLE. The new user can subsequently be set as the owner of a database via 'dbaas postgres database create --owner'.

Required values to run command:

* Cluster Id
* User (name)
* Password`,
			Example:   `ionosctl dbaas postgres user create --cluster-id CLUSTER_ID --user appuser --password 'S3cr3tPassw0rd'`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster the user is created in (must be AVAILABLE)")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "Name of the PostgreSQL login role to create. Must not collide with a reserved system name (e.g. postgres)")
	c.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "Login password for the new user. Minimum 10 characters")

	return c
}

func preRunCreateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.ArgUser, constants.ArgPassword)
}

func runCreateCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	username := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	password := viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))

	user, _, err := client.Must().PostgresClient.UsersApi.UsersPost(
		context.Background(),
		clusterId,
	).User(
		psql.User{
			Properties: psql.UserProperties{
				Username: username,
				Password: &password,
			},
		},
	).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(user)
}
