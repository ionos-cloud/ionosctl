package vulnerabilities

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "DataSource", Default: true, Format: func(item map[string]any) any {
		props := table.Navigate(item, "properties.dataSource")
		if props == nil {
			return nil
		}
		m, _ := props.(map[string]any)
		if m == nil {
			return nil
		}
		id, _ := m["id"].(string)
		url, _ := m["url"].(string)
		if id != "" && url != "" {
			return fmt.Sprintf("%s (%s)", id, url)
		}
		return nil
	}},
	{Name: "Score", JSONPath: "properties.score", Default: true},
	{Name: "Severity", JSONPath: "properties.severity", Default: true},
	{Name: "Fixable", JSONPath: "properties.fixable", Default: true},
	{Name: "PublishedAt", JSONPath: "metadata.publishedAt", Default: true},
	{Name: "UpdatedAt", JSONPath: "metadata.updatedAt"},
	{Name: "Affects", Format: func(item map[string]any) any {
		affects := table.Navigate(item, "properties.affects")
		if affects == nil {
			return nil
		}
		arr, _ := affects.([]any)
		if arr == nil {
			return nil
		}
		var parts []string
		for _, a := range arr {
			m, _ := a.(map[string]any)
			if m == nil {
				continue
			}
			name, _ := m["name"].(string)
			version, _ := m["version"].(string)
			parts = append(parts, fmt.Sprintf("%s (%s)", name, version))
		}
		return strings.Join(parts, ", ")
	}},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "Recommendations", JSONPath: "properties.recommendations"},
	{Name: "References", JSONPath: "properties.references"},
	{Name: "Href", JSONPath: "href"},
}

func VulnerabilitiesCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vulnerabilities",
			Aliases:          []string{"v", "vuln", "vulnerability"},
			Short:            "Vulnerabilities Operations",
			Long:             "Manage container registry vulnerabilities.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(VulnerabilitiesGetCmd())
	cmd.AddCommand(VulnerabilitiesListCmd())

	return cmd
}
