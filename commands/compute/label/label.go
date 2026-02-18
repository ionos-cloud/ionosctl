package label

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

var (
	defaultLabelCols = []string{"URN", "Key", "Value", "ResourceType", "ResourceId"}
)

func LabelCmd() *core.Command {
	labelCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "label",
			Short:            "Label Operations",
			Long:             "The sub-commands of `ionosctl label` allow you to get, list, add, remove Labels from a Resource.",
			TraverseChildren: true,
		},
	}
	globalFlags := labelCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLabelCols, tabheaders.ColsMessage(defaultLabelCols))
	_ = viper.BindPFlag(core.GetFlagName(labelCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultLabelCols, cobra.ShellCompDirectiveNoFileComp
	})

	labelCmd.AddCommand(LabelListCmd())
	labelCmd.AddCommand(LabelGetCmd())
	labelCmd.AddCommand(LabelGetByUrnCmd())
	labelCmd.AddCommand(LabelAddCmd())
	labelCmd.AddCommand(LabelRemoveCmd())

	return core.WithConfigOverride(labelCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
