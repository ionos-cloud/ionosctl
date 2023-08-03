package templates

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
			req := client.Must().MongoClient.TemplatesApi.TemplatesGet(context.Background())

			if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
				req = req.Limit(viper.GetInt32(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
				req = req.Offset(viper.GetInt32(f))
			}

			ls, _, err := req.Execute()
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
