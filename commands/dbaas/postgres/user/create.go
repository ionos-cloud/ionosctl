package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
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
			LongDesc:  `Create a new user in the given cluster`,
			Example:   `ionosctl dbaas postgres user create --cluster-id <cluster-id> --user <user> --password <password>`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The ID of the Postgres cluster")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "The name of the user")
	c.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "The password for the user")

	return c
}

func preRunCreateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.ArgUser, constants.ArgPassword)
}

func runCreateCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	username := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	password := viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))

	userProps := psql.UserProperties{Username: &username, Password: &password}
	user, _, err := client.Must().PostgresClient.UsersApi.UsersPost(
		context.Background(),
		clusterId,
	).User(psql.User{Properties: &userProps}).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasPostgresUser, user, defaultCols)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
