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
)

func UpdateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "update",
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Update a user",
			LongDesc:  `Update the specified user from the given cluster. Only changing their password is supported. System users cannot be patched.`,
			Example:   `ionosctl dbaas postgres user update --cluster-id <cluster-id> --user <user> --password <password>`,
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
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
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgUser,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "The password for the user")

	return c
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.ArgUser, constants.ArgPassword)
}

func runUpdateCmd(c *core.CommandConfig) error {
	clusterId, err := c.Command.Command.Flags().GetString(constants.FlagClusterId)
	if err != nil {
		return err
	}

	username, err := c.Command.Command.Flags().GetString(constants.ArgUser)
	if err != nil {
		return err
	}

	password, err := c.Command.Command.Flags().GetString(constants.ArgPassword)
	if err != nil {
		return err
	}

	user, _, err := client.Must().PostgresClient.UsersApi.UsersPatch(
		context.Background(),
		clusterId,
		username,
	).UsersPatchRequest(psql.UsersPatchRequest{Properties: psql.PatchUserProperties{Password: &password}}).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasPostgresUser, user, defaultCols)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
