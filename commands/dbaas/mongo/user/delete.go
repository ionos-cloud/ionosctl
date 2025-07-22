package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
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

			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete user %s", user),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf("operation canceled by confirmation check")
			}

			u, _, err := client.Must().MongoClient.UsersApi.
				ClustersUsersDelete(context.Background(), clusterId, user).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			uConverted, err := resource2table.ConvertDbaasMongoUserToTable(u)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(u, uConverted, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
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

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig, clusterId string) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all users"))
	xs, _, err := client.Must().MongoClient.UsersApi.ClustersUsersGet(c.Context, clusterId).Execute()
	if err != nil {
		return err
	}

	return functional.ApplyAndAggregateErrors(xs.GetItems(), func(x sdkgo.User) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete user %s", x.Properties.Username), viper.GetBool(constants.ArgForce))
		if !yes {
			return fmt.Errorf("user %s skipped by confirmation check", x.Properties.Username)
		}
		_, _, delErr := client.Must().MongoClient.UsersApi.ClustersUsersDelete(c.Context, clusterId, x.Properties.Username).Execute()
		if delErr != nil {
			return fmt.Errorf("failed deleting one of the resources: %w", delErr)
		}
		return nil
	})
}
