package contract

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ContractGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "contract",
		Resource:   "contract",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get information about the Contract Resources on your account",
		LongDesc:   "Use this command to get information about the Contract Resources on your account. Use `--resource-limits` flag to see specific Contract Resources Limits.",
		Example:    `ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunContractGet,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgResourceLimits, "", "", "Specify Resource Limits to see details about it")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceLimits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CORES", "RAM", "HDD", "SSD", "DAS", "IPS", "K8S", "NLB", "NAT"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
