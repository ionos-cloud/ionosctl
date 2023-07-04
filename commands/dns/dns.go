package dns

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/record"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

func DNSCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dns",
			Short:            "DNS API",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(zone.ZoneCommand())
	cmd.AddCommand(record.RecordCommand())

	return cmd
}
