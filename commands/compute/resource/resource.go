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
			Use:     "resource",
			Aliases: []string{"res"},
			Short:   "Read-only overview of shareable cloud resources",
			Long: `A read-only overview of the cloud resources on your contract that can be shared with Groups - datacenters, snapshots, images, IP blocks, PCCs, backup units and Kubernetes clusters. Each entry shows the resource's ID, name, type, provisioning state, and whether it is protected by secure (two-factor) authentication.

Use these commands to discover resource IDs and types before granting a Group access to them with ` + "`ionosctl compute share create`" + `. This tree is read-only: it does not create or modify the resources themselves (manage those with their own command trees, e.g. ` + "`ionosctl compute datacenter`" + `).`,
			TraverseChildren: true,
		},
	}
	resourceCmd.AddColsFlag(allResourceCols)

	resourceCmd.AddCommand(ResourceListCmd())
	resourceCmd.AddCommand(ResourceGetCmd())

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
