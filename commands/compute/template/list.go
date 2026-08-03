package template

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func TemplateListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "template",
		Resource:  "template",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List the CUBE templates and their fixed cores/RAM/storage",
		LongDesc:  "Use this command to list the predefined CUBE templates. Each row shows the fixed cores, RAM and storage size (and any GPUs) of a template; the TemplateId is what you pass as `--template-id` when creating a CUBE server.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.TemplatesFiltersUsage(),
		Example: `# List all CUBE templates
ionosctl compute template list

# Find templates by name
ionosctl compute template list --filters name=CUBE`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTemplateList,
		InitClient: true,
	})

	return cmd
}
