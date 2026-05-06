package label

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

var allLabelCols = []table.Column{
	{Name: "URN", JSONPath: "id", Default: true},
	{Name: "Key", JSONPath: "properties.key", Default: true},
	{Name: "Value", JSONPath: "properties.value", Default: true},
	{Name: "ResourceType", JSONPath: "properties.resourceType", Default: true},
	{Name: "ResourceId", JSONPath: "properties.resourceId", Default: true},
}

func LabelCmd() *core.Command {
	labelCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "label",
			Short:            "Label Operations",
			Long:             "The sub-commands of `ionosctl compute label` allow you to get, list, add, remove Labels from a Resource.",
			TraverseChildren: true,
		},
	}
	labelCmd.AddColsFlag(allLabelCols)

	labelCmd.AddCommand(LabelListCmd())
	labelCmd.AddCommand(LabelGetCmd())
	labelCmd.AddCommand(LabelGetByUrnCmd())
	labelCmd.AddCommand(LabelAddCmd())
	labelCmd.AddCommand(LabelRemoveCmd())

	return core.WithConfigOverride(labelCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
