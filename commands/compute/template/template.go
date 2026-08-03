package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allTemplateCols = []table.Column{
	{Name: "TemplateId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Cores", JSONPath: "properties.cores", Default: true},
	{Name: "RAM", JSONPath: "properties.ram", Default: true},
	{Name: "StorageSize", JSONPath: "properties.storageSize", Default: true},
	{Name: "GPUs", Default: true, Format: func(item map[string]any) any {
		gpus := table.Navigate(item, "properties.gpus")
		if gpus == nil {
			return ""
		}
		gpuSlice, ok := gpus.([]any)
		if !ok || len(gpuSlice) == 0 {
			return ""
		}
		formatted := make([]string, 0, len(gpuSlice))
		for _, g := range gpuSlice {
			gMap, ok := g.(map[string]any)
			if !ok {
				continue
			}
			model, _ := gMap["model"].(string)
			count, _ := gMap["count"].(float64)
			if model != "" && count > 0 {
				formatted = append(formatted, fmt.Sprintf("%dx %s", int(count), model))
			}
		}
		sort.Strings(formatted)
		return strings.Join(formatted, ", ")
	}},
}

func TemplateCmd() *core.Command {
	templateCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "template",
			Aliases: []string{"tpl"},
			Short:   "Browse the predefined CUBE server templates",
			Long: `The ` + "`ionosctl compute template`" + ` command lists the predefined templates used to create CUBE servers.

A CUBE server is not sized freely: instead you pick one template, which is a fixed bundle of cores, RAM and a fixed-size NVMe Direct Attached Storage (DAS) volume (plus any bundled GPUs). You reference the chosen template by its ID via ` + "`--template-id`" + ` when running ` + "`ionosctl compute server create --type CUBE`" + `. Templates do not apply to ENTERPRISE or VCPU servers, whose cores/RAM/storage are specified individually.

Use ` + "`list`" + ` to see all templates and ` + "`get`" + ` to inspect one. This command is read-only.`,
			TraverseChildren: true,
		},
	}
	templateCmd.AddColsFlag(allTemplateCols)

	templateCmd.AddCommand(TemplateListCmd())
	templateCmd.AddCommand(TemplateGetCmd())

	return core.WithConfigOverride(templateCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
