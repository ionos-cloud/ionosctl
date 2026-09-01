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
			Use:     "ipfailover",
			Aliases: []string{"ipf"},
			Short:   "IP Failover Operations",
			Long: `The sub-commands of ` + "`" + `ionosctl compute ipfailover` + "`" + ` let you manage IP failover groups on a LAN.

An IP failover group gives you a single reserved public IP that "floats" between the NICs of several servers on the same LAN, providing high availability. At any moment one NIC is the active master that actually holds the IP; the others are standby members. If the master's server fails, the IP is reassigned to a standby member so the service keeps answering on the same address.

How it fits together:
  * The floating IP must come from a reserved IP block in the SAME region/location as the datacenter (see ` + "`" + `ionosctl ipblock` + "`" + `). Ad-hoc/DHCP IPs cannot be used.
  * Every participating NIC must be on the same LAN, and ideally on servers spread across Availability Zones for real redundancy.
  * A NIC can hold the primary failover IP for only one failover group at a time.
  * Actual failover detection/switchover is driven at the OS level (e.g. keepalived/heartbeat) on the servers; this command only registers which IP/NIC pairs belong to the group.

Use ` + "`" + `add` + "`" + ` to register the master NIC (and repeat for each member using the same IP), ` + "`" + `list` + "`" + ` to view the group on a LAN, and ` + "`" + `remove` + "`" + ` to take a NIC out of the group.`,
			TraverseChildren: true,
		},
	}
	ipfailoverCmd.AddColsFlag(allIpFailoverCols)

	ipfailoverCmd.AddCommand(IpFailoverListCmd())
	ipfailoverCmd.AddCommand(IpFailoverAddCmd())
	ipfailoverCmd.AddCommand(IpFailoverRemoveCmd())

	return core.WithConfigOverride(ipfailoverCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
