package file

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "file",
			Aliases: []string{"f"},
			Short:   "Import/export a zone as a BIND file",
			Long: `Import or export a whole zone as a single BIND-format file (RFC 1035) — the same text format used by standard name servers.

'get' exports every record in the zone as text (handy for backup or migration); 'update' replaces the zone's records with the contents of a BIND file, so you can edit records in bulk instead of one 'dns record' call at a time.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(getCmd())
	cmd.AddCommand(updateCmd())

	return cmd
}
