package user

import (
	"context"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a MongoDB user",
		Example:   "ionosctl dbaas mongo user get",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			database := viper.GetString(core.GetFlagName(c.NS, flagDatabase))
			user := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
			c.Printer.Verbose("Getting User by ID %s...")
			u, _, err := c.DbaasMongoServices.Users().Get(clusterId, database, user)
			if err != nil {
				return err
			}
			return c.Printer.Print(getUserPrint(c, &[]sdkgo.User{u}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagDatabase, "d", "", "The authentication database")
	cmd.AddStringFlag(constants.ArgUser, "u", "", "The authentication username")
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
