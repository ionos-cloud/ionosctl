package user

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/internal/confirm"
	"github.com/ionos-cloud/ionosctl/internal/functional"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "delete",
		Aliases:   []string{"g"},
		ShortDesc: "Delete a MongoDB user",
		Example:   "ionosctl dbaas mongo user delete",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagClusterId, constants.ArgAll}, []string{constants.FlagClusterId, constants.FlagName})
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c, clusterId)
			}
			user := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
			u, _, err := c.DbaasMongoServices.Users().Delete(clusterId, user)
			if err != nil {
				return err
			}
			return c.Printer.Print(getUserPrint(c, &[]sdkgo.User{u}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(FlagDatabase, FlagDatabaseShort, "", "The authentication database")
	cmd.AddStringFlag(constants.ArgUser, "", "", "The authentication username")
	cmd.AddBoolFlag(constants.ArgAll, "a", false, "Delete all users in a cluster")

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig, clusterId string) error {
	client, err := config.GetClient()
	if err != nil {
		return err
	}
	c.Printer.Verbose("Deleting all users")
	xs, _, err := client.MongoClient.UsersApi.ClustersUsersGet(c.Context, clusterId).Execute()
	if err != nil {
		return err
	}

	return functional.ApplyOrFail(*xs.GetItems(), func(x sdkgo.User) error {
		yes := confirm.Ask(fmt.Sprintf("delete user %s", *x.Properties.Username), viper.GetBool(constants.ArgForce))
		if yes {
			_, _, delErr := client.MongoClient.UsersApi.ClustersUsersDelete(c.Context, clusterId, *x.Properties.Username).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting one of the resources: %w", delErr)
			}
		}
		return nil
	})
}
