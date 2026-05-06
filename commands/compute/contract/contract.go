package contract

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	allContractCols = []table.Column{
		{Name: "ContractNumber", JSONPath: "properties.contractNumber", Default: true},
		{Name: "Owner", JSONPath: "properties.owner", Default: true},
		{Name: "Status", JSONPath: "properties.status", Default: true},
		{Name: "RegistrationDomain", JSONPath: "properties.regDomain", Default: true},
		{Name: "CoresPerServer", JSONPath: "properties.resourceLimits.coresPerServer"},
		{Name: "CoresPerContract", JSONPath: "properties.resourceLimits.coresPerContract"},
		{Name: "CoresProvisioned", JSONPath: "properties.resourceLimits.coresProvisioned"},
		{Name: "RamPerServer", JSONPath: "properties.resourceLimits.ramPerServer"},
		{Name: "RamPerContract", JSONPath: "properties.resourceLimits.ramPerContract"},
		{Name: "RamProvisioned", JSONPath: "properties.resourceLimits.ramProvisioned"},
		{Name: "HddLimitPerVolume", JSONPath: "properties.resourceLimits.hddLimitPerVolume"},
		{Name: "HddLimitPerContract", JSONPath: "properties.resourceLimits.hddLimitPerContract"},
		{Name: "HddVolumeProvisioned", JSONPath: "properties.resourceLimits.hddVolumeProvisioned"},
		{Name: "SsdLimitPerVolume", JSONPath: "properties.resourceLimits.ssdLimitPerVolume"},
		{Name: "SsdLimitPerContract", JSONPath: "properties.resourceLimits.ssdLimitPerContract"},
		{Name: "SsdVolumeProvisioned", JSONPath: "properties.resourceLimits.ssdVolumeProvisioned"},
		{Name: "DasVolumeProvisioned", JSONPath: "properties.resourceLimits.dasVolumeProvisioned"},
		{Name: "ReservableIps", JSONPath: "properties.resourceLimits.reservableIps"},
		{Name: "ReservedIpsOnContract", JSONPath: "properties.resourceLimits.reservedIpsOnContract"},
		{Name: "ReservedIpsInUse", JSONPath: "properties.resourceLimits.reserverIpsInUse"},
		{Name: "K8sClusterLimitTotal", JSONPath: "k8sClusterLimitTotal"},
		{Name: "K8sClustersProvisioned", JSONPath: "k8sClustersProvisioned"},
		{Name: "NlbLimitTotal", JSONPath: "properties.resourceLimits.nlbLimitTotal"},
		{Name: "NlbProvisioned", JSONPath: "properties.resourceLimits.nlbProvisioned"},
		{Name: "NatGatewayLimitTotal", JSONPath: "properties.resourceLimits.natGatewayLimitTotal"},
		{Name: "NatGatewayProvisioned", JSONPath: "properties.resourceLimits.natGatewayProvisioned"},
	}

	contractCoresCols = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "CoresPerServer", "CoresPerContract", "CoresProvisioned"}
	contractRamCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "RamPerServer", "RamPerContract", "RamProvisioned"}
	contractHddCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "HddLimitPerVolume", "HddLimitPerContract", "HddVolumeProvisioned"}
	contractSsdCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "SsdLimitPerVolume", "SsdLimitPerContract", "SsdVolumeProvisioned"}
	contractDasCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "DasVolumeProvisioned"}
	contractIpsCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "ReservableIps", "ReservedIpsOnContract", "ReservedIpsInUse"}
	contractK8sCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "K8sClusterLimitTotal", "K8sClustersProvisioned"}
	contractNatCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "NatGatewayLimitTotal", "NatGatewayProvisioned"}
	contractNlbCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "NlbLimitTotal", "NlbProvisioned"}
)

func ContractCmd() *core.Command {
	contractCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "contract",
			Aliases:          []string{"c"},
			Short:            "Contract Resources Operations",
			Long:             "The sub-command of `ionosctl compute contract` allows you to see information about Contract Resources for your account.",
			TraverseChildren: true,
		},
	}
	contractCmd.AddColsFlag(allContractCols)

	contractCmd.AddCommand(ContractGetCmd())

	return core.WithConfigOverride(contractCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
