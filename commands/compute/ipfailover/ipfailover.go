package ipfailover

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allIpFailoverCols = []table.Column{
	{Name: "NicId", JSONPath: "nicUuid", Default: true},
	{Name: "Ip", JSONPath: "ip", Default: true},
}

func IpfailoverCmd() *core.Command {
	ipfailoverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipfailover",
			Aliases:          []string{"ipf"},
			Short:            "IP Failover Operations",
			Long:             "The sub-command of `ionosctl compute ipfailover` allows you to see information about IP Failovers groups available on a LAN, to add/remove IP Failover group from a LAN.",
			TraverseChildren: true,
		},
	}
	ipfailoverCmd.AddColsFlag(allIpFailoverCols)

	ipfailoverCmd.AddCommand(IpFailoverListCmd())
	ipfailoverCmd.AddCommand(IpFailoverAddCmd())
	ipfailoverCmd.AddCommand(IpFailoverRemoveCmd())

	return core.WithConfigOverride(ipfailoverCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
