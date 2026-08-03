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
			Use:     "location",
			Aliases: []string{"loc"},
			Short:   "Discover the IONOS regions where you can provision resources",
			Long: `The ` + "`ionosctl compute location`" + ` command lists the physical locations (data centers) where IONOS Cloud resources can be created.

A location ID has the form ` + "`region/location`" + `, where the region is a two-letter country code and the location is a site within it, e.g. ` + "`de/txl`" + ` (Berlin), ` + "`de/fra`" + ` (Frankfurt), ` + "`us/las`" + ` (Las Vegas). You pass this ID as the ` + "`--location`" + ` flag when creating datacenters and other resources.

Each location advertises the features it supports (e.g. flow logs), available image aliases, and its CPU architectures. Use ` + "`ionosctl compute location cpu`" + ` to drill into the CPU families (vendor, max cores, max RAM, enabled features) available at a specific site, which determines the ` + "`--cpu-family`" + ` values you can request when creating servers. This command is read-only.`,
			TraverseChildren: true,
		},
	}
	locationCmd.AddColsFlag(allLocationCols)

	locationCmd.AddCommand(LocationListCmd())
	locationCmd.AddCommand(LocationGetCmd())
	locationCmd.AddCommand(CpuCmd())

	return core.WithConfigOverride(locationCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
