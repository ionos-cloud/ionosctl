package distribution

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/distribution/routingrules"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Domain", JSONPath: "properties.domain", Default: true},
	{Name: "CertificateId", JSONPath: "properties.certificateId", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "distribution",
			Short: "Manage CDN distributions: the domain, certificate, and routing rules that serve your content from the edge",
			Long: `A CDN distribution is the top-level CDN resource. It maps a public DOMAIN (and an optional HTTPS CERTIFICATE) to a set of ROUTING RULES, where each rule matches requests by URL path prefix and forwards them to an upstream origin host.

Use these sub-commands to list and inspect distributions, create and update them (domain, certificate binding, and routing rules), delete them, and view the effective routing rules of a distribution (see 'routingrules').

Distribution state is reported in the State column: AVAILABLE (serving), BUSY (a change is being provisioned), FAILED, or UNKNOWN.`,
			Aliases:          []string{"ds"},
			TraverseChildren: true,
		},
	}
	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(FindByID())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Create())
	cmd.AddCommand(Update())
	cmd.AddCommand(routingrules.Root())
	return cmd
}
