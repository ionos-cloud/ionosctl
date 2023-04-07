package user

import (
	"context"
	"fmt"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
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
			user := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
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
	cmd.AddStringFlag(constants.FlagName, "", "", "The authentication username")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all users in a cluster")
	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip y/n checks")

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig, clusterId string) error {
	client, err := client2.Get()
	if err != nil {
		return err
	}
	c.Printer.Verbose("Deleting all users")
	xs, _, err := client.MongoClient.UsersApi.ClustersUsersGet(c.Context, clusterId).Execute()
	if err != nil {
		return err
	}

	return functional.ApplyAndAggregateErrors(*xs.GetItems(), func(x sdkgo.User) error {
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
