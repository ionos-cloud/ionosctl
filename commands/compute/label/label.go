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
			Use:   "label",
			Short: "Label Operations",
			Long: `A Label is a key/value pair attached to a cloud resource to help you organize, group and filter your infrastructure (for example env=prod, team=payments, cost-center=1234).

Domain model:
  * A label lives ON a resource. It is identified by a resource + a key; a resource can hold at most one label per key (adding an existing key overwrites its value).
  * Labels can be attached to several compute resource types: datacenter, server, volume, snapshot, ipblock and (private) image. You select which with --resource-type plus that resource's id flag(s).
  * Nested resources need their parent(s): a server label needs --datacenter-id and --server-id; a volume label needs --datacenter-id and --volume-id; a datacenter label needs only --datacenter-id. Snapshots, ipblocks and images are contract-level, so only their own id is needed.
  * Every label also has a globally unique URN of the form urn:label:<resource_type>:<resource_id>:<key>. Use ` + "`label get-by-urn`" + ` to fetch a label directly by URN, and ` + "`label list`" + ` (without --resource-type) to list labels across all resources.

The sub-commands of ` + "`ionosctl compute label`" + ` let you list, get (by resource or by URN), add and remove Labels.`,
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
