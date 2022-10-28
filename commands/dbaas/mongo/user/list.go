package user

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
)

func UserListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "Retrieves a list of MongoDB users.",
		Example:   "ionosctl dbaas mongo user list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagRequired("")
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Templates...")
			ls, r, err := c.DbaasMongoServices.Users().List("")
			if err != nil {
				return err
			}
			return c.Printer.Print(getUserPrint(r, c, ls.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(config.ArgCols, "", allCols, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
