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
			Use:     "contract",
			Aliases: []string{"c"},
			Short:   "View your contract and account resource limits",
			Long: `The ` + "`ionosctl compute contract`" + ` command shows your IONOS Cloud contract: the contract number, the account owner, the registration (billing) status, and the resource limits (quotas) attached to your contract.

Resource limits are hard caps enforced by the platform. Each limit is reported both as a maximum and as the amount currently provisioned, so you can tell how much headroom you have before a provisioning request is rejected. Limits are grouped by resource type:

* Cores: vCPUs per single server and across the whole contract
* RAM: memory per single server and across the whole contract (in MB)
* HDD / SSD / DAS: block-storage capacity per volume and per contract (in GB)
* IPS: reservable IPv4 addresses, and how many are reserved / in use
* K8S: total Managed Kubernetes clusters allowed
* NLB / NAT: Network Load Balancers and NAT Gateways allowed

To raise a limit, contact IONOS support. This command is read-only.`,
			TraverseChildren: true,
		},
	}
	contractCmd.AddColsFlag(allContractCols)

	contractCmd.AddCommand(ContractGetCmd())

	return core.WithConfigOverride(contractCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
