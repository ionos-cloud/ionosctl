package templates

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
)

func TemplatesListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "templates",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Mongo Templates",
		LongDesc:  "Retrieves a list of valid templates. These templates can be used to create MongoDB clusters; they contain properties, such as number of cores, RAM, and the storage size.",
		Example:   "ionosctl dbaas mongo templates list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Templates...")
			ls, r, err := c.DbaasMongoServices.Templates().List()
			if err != nil {
				return err
			}
			return c.Printer.Print(getTemplatesPrint(r, c, ls.GetItems()))
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
