package template

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var templateCols = []table.Column{
	{Name: "TemplateId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Edition", JSONPath: "properties.edition", Default: true},
	{Name: "Cores", JSONPath: "properties.cores", Default: true},
	{Name: "Ram", JSONPath: "properties.ram", Default: true},
	{Name: "StorageSize", JSONPath: "properties.storageSize", Default: true},
}

func TemplateCmd() *core.Command {
	templateCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "template",
			Aliases: []string{"templates", "t"},
			Short:   "MongoDB Instance Template Operations",
			Long: "The sub-commands of `ionosctl dbaas mongo-v2 template` allow you to view MongoDB Instance Templates.\n\n" +
				"An instance template is a predefined bundle of cores, RAM (GB) and storage size (GB) per instance. BUSINESS clusters are sized by one, passed as `--template-id` to `cluster create`.",
			TraverseChildren: true,
		},
	}

	templateCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(templateCols))
	_ = templateCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(templateCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	templateCmd.AddCommand(TemplateListCmd())
	templateCmd.AddCommand(TemplateGetCmd())

	return templateCmd
}
