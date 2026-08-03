package contract

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ContractGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "contract",
		Resource:  "contract",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get your contract number, owner, status and resource limits",
		LongDesc: `Use this command to view your contract details and the resource limits (quotas) enforced on your account.

By default all limits are shown. Use ` + "`--resource-limits`" + ` to focus the output on a single resource group; each group also shows the amount already provisioned so you can gauge remaining headroom.`,
		Example: `# Show full contract info and all resource limits
ionosctl compute contract get

# Focus on core (vCPU) limits only
ionosctl compute contract get --resource-limits CORES`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunContractGet,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgResourceLimits, "", "", "Restrict the output to one resource-limit group. One of: CORES (vCPUs), RAM, HDD, SSD, DAS (block storage), IPS (reservable IPs), K8S (Kubernetes clusters), NLB (Network Load Balancers), NAT (NAT Gateways)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceLimits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CORES", "RAM", "HDD", "SSD", "DAS", "IPS", "K8S", "NLB", "NAT"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
