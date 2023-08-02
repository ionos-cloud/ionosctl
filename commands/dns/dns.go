package dns

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/record"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

const (
	DefaultApiURL = "dns.de-fra.ionos.com"
)

func DNSCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dns",
			Short:            "The sub-commands of `ionosctl dns` allow you to manage your DNS Zones (domains) and Records (where the traffic should be redirected)",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(zone.ZoneCommand())
	cmd.AddCommand(record.RecordCommand())

	return cmd
}
