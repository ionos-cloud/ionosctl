package contract

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultContractCols = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain"}
	allContractCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "CoresPerServer", "CoresPerContract", "CoresProvisioned", "RamPerServer", "RamPerContract", "RamProvisioned",
		"HddLimitPerVolume", "HddLimitPerContract", "HddVolumeProvisioned", "SsdLimitPerVolume", "SsdLimitPerContract", "SsdVolumeProvisioned", "DasVolumeProvisioned", "ReservableIps", "ReservedIpsOnContract",
		"ReservedIpsInUse", "K8sClusterLimitTotal", "K8sClustersProvisioned", "NlbLimitTotal", "NlbProvisioned", "NatGatewayLimitTotal", "NatGatewayProvisioned"}
)

func ContractCmd() *core.Command {
	contractCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "contract",
			Aliases:          []string{"c"},
			Short:            "Contract Resources Operations",
			Long:             "The sub-command of `ionosctl contract` allows you to see information about Contract Resources for your account.",
			TraverseChildren: true,
		},
	}
	globalFlags := contractCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultContractCols, tabheaders.ColsMessage(allContractCols))
	_ = viper.BindPFlag(core.GetFlagName(contractCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = contractCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allContractCols, cobra.ShellCompDirectiveNoFileComp
	})

	contractCmd.AddCommand(ContractGetCmd())

	return core.WithConfigOverride(contractCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
