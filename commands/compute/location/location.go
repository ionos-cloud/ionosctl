package location

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allLocationCols = []table.Column{
	{Name: "LocationId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Features", JSONPath: "properties.features"},
	{Name: "ImageAliases", JSONPath: "properties.imageAliases"},
	{Name: "CpuFamily", JSONPath: "properties.cpuArchitecture.*.cpuFamily", Default: true},
	{Name: "CpuEnabledFeatures", JSONPath: "properties.cpuArchitecture.*.enabledFeatures"},
}

func LocationCmd() *core.Command {
	locationCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "location",
			Aliases:          []string{"loc"},
			Short:            "Location Operations",
			Long:             "The sub-command of `ionosctl compute location` allows you to see information about locations available to create objects.",
			TraverseChildren: true,
		},
	}
	locationCmd.AddColsFlag(allLocationCols)

	locationCmd.AddCommand(LocationListCmd())
	locationCmd.AddCommand(LocationGetCmd())
	locationCmd.AddCommand(CpuCmd())

	return core.WithConfigOverride(locationCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
