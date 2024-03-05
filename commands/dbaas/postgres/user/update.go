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
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UpdateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "update",
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Update a user",
			LongDesc:  "Update a new user in the given database cluster",
			Example:
				"ionosctl dbaas postgres user update --cluster-id <cluster-id> --user <user> --password <password>",
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveDefault
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "The name of the new user")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgUser,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveDefault
		},
	)

	c.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "The password of the new user")

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
	).UsersPatchRequest(ionoscloud.UsersPatchRequest{Properties: &ionoscloud.PatchUserProperties{Password: &password}}).Execute()
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
