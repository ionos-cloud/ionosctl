package template

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultTemplateCols = []string{"TemplateId", "Name", "Cores", "RAM", "StorageSize", "GPUs"}
)

func TemplateCmd() *core.Command {
	templateCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "template",
			Aliases:          []string{"tpl"},
			Short:            "Template Operations",
			Long:             "The sub-commands of `ionosctl template` allow you to see information about the Templates available.",
			TraverseChildren: true,
		},
	}
	globalFlags := templateCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultTemplateCols, tabheaders.ColsMessage(defaultTemplateCols))
	_ = viper.BindPFlag(core.GetFlagName(templateCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = templateCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTemplateCols, cobra.ShellCompDirectiveNoFileComp
	})

	templateCmd.AddCommand(TemplateListCmd())
	templateCmd.AddCommand(TemplateGetCmd())

	return core.WithConfigOverride(templateCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
