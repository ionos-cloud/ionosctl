package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Get user",
			LongDesc:  `Retrieve a single user of a cluster by name. The response does not include the password.` + "\n\n" + `Required values to run command:` + "\n\n" + `* Cluster Id` + "\n" + `* User (name)`,
			Example:   `ionosctl dbaas postgres user get --cluster-id CLUSTER_ID --user appuser`,
			PreCmdRun: preRunGetCmd,
			CmdRun:    runGetCmd,
		},
	)
	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster the user belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "Name of the user to retrieve")
	_ = c.Command.RegisterFlagCompletionFunc(constants.ArgUser, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UserNames(c), cobra.ShellCompDirectiveNoFileComp
	})

	return c
}

func preRunGetCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.ArgUser)
}

func runGetCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	username := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))

	user, _, err := client.Must().PostgresClient.UsersApi.UsersGet(context.Background(), clusterId, username).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(user)
}
