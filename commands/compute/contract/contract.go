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
	contractCoresCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "CoresPerServer", "CoresPerContract", "CoresProvisioned"}
	contractRamCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "RamPerServer", "RamPerContract", "RamProvisioned"}
	contractHddCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "HddLimitPerVolume", "HddLimitPerContract", "HddVolumeProvisioned"}
	contractSsdCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "SsdLimitPerVolume", "SsdLimitPerContract", "SsdVolumeProvisioned"}
	contractDasCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "DasVolumeProvisioned"}
	contractIpsCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "ReservableIps", "ReservedIpsOnContract", "ReservedIpsInUse"}
	contractK8sCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "K8sClusterLimitTotal", "K8sClustersProvisioned"}
	contractNatCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "NatGatewayLimitTotal", "NatGatewayProvisioned"}
	contractNlbCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "NlbLimitTotal", "NlbProvisioned"}
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
