package vulnerabilities

import (
	"fmt"
	"strings"

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
			Use:     "vulnerabilities",
			Aliases: []string{"v", "vuln", "vulnerability"},
			Short:   "Inspect vulnerability scan findings for artifacts",
			Long: `A vulnerability is a security finding produced by scanning a pushed artifact against public CVE data sources. Findings are only available when the registry's vulnerabilityScanning feature is enabled (a paid add-on; see 'container-registry registry create --vuln-scan').

Each finding carries a CVSS score, a severity (unknown, none, info, low, medium, high, critical), whether a fix is available (fixable), the affected packages/versions, and remediation recommendations. These commands are read-only.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(VulnerabilitiesGetCmd())
	cmd.AddCommand(VulnerabilitiesListCmd())

	return cmd
}
