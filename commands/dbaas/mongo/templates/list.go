package templates

import (
	"context"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
)

func TemplatesListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
			var limitPtr *int32 = nil
			if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
				limit := viper.GetInt32(f)
				limitPtr = &limit
			}
			var offsetPtr *int32 = nil
			if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
				offset := viper.GetInt32(f)
				offsetPtr = &offset
			}
			ls, _, err := c.DbaasMongoServices.Templates().List(limitPtr, offsetPtr)
			if err != nil {
				return err
			}
			return c.Printer.Print(getTemplatesPrint(c, ls.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", allCols, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	return cmd
}
