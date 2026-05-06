package resource

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allResourceCols = []table.Column{
	{Name: "ResourceId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "SecAuthProtection", JSONPath: "properties.secAuthProtection", Default: true},
	{Name: "Type", JSONPath: "type", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ResourceCmd() *core.Command {
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "resource",
			Aliases:          []string{"res"},
			Short:            "Resource Operations",
			Long:             "The sub-commands of `ionosctl compute resource` allow you to list, get Resources.",
			TraverseChildren: true,
		},
	}
	resourceCmd.AddColsFlag(allResourceCols)

	resourceCmd.AddCommand(ResourceListCmd())
	resourceCmd.AddCommand(ResourceGetCmd())

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
