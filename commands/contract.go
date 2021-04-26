package commands

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func contract() *builder.Command {
	ctx := context.TODO()
	contractCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "contract",
			Short:            "Contract Resources Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl contract` + "`" + ` allows you to see information about Contract Resources for your account.`,
			TraverseChildren: true,
		},
	}
	globalFlags := contractCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultContractCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(contractCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	get := builder.NewCommand(ctx, contractCmd, noPreRun, RunContractGet, "get", "Get information about the Contract Resources on your account",
		"Use this command to get information about the Contract Resources on your account. Use `--resource-limits` flag to see specific Contract Resources Limits.",
		getContractExample, true)
	get.AddStringFlag(config.ArgResourceLimits, "", "", "Specify Resource Limits to see details about it")
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgResourceLimits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CORES", "RAM", "HDD", "SSD", "IPS", "K8S"}, cobra.ShellCompDirectiveNoFileComp
	})

	return contractCmd
}

func RunContractGet(c *builder.CommandConfig) error {
	contractResource, _, err := c.Contracts().Get()
	if err != nil {
		return err
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceLimits)) {
		switch strings.ToUpper(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceLimits))) {
		case "CORES":
			return c.Printer.Print(getContractPrint(c, getContract(&contractResource), contractCoresCols))
		case "RAM":
			return c.Printer.Print(getContractPrint(c, getContract(&contractResource), contractRamCols))
		case "HDD":
			return c.Printer.Print(getContractPrint(c, getContract(&contractResource), contractHddCols))
		case "SSD":
			return c.Printer.Print(getContractPrint(c, getContract(&contractResource), contractSsdCols))
		case "IPS":
			return c.Printer.Print(getContractPrint(c, getContract(&contractResource), contractIpsCols))
		case "K8S":
			return c.Printer.Print(getContractPrint(c, getContract(&contractResource), contractK8sCols))
		}
	}
	return c.Printer.Print(getContractPrint(c, getContract(&contractResource), getContractCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())))
}

// Output Printing

var (
	defaultContractCols = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain"}
	contractCoresCols   = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "CoresPerServer", "CoresPerContract", "CoresProvisioned"}
	contractRamCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "RamPerServer", "RamPerContract", "RamProvisioned"}
	contractHddCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "HddLimitPerVolume", "HddLimitPerContract", "HddVolumeProvisioned"}
	contractSsdCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "SsdLimitPerVolume", "SsdLimitPerContract", "SsdVolumeProvisioned"}
	contractIpsCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "ReservableIps", "ReservedIpsOnContract", "ReservedIpsInUse"}
	contractK8sCols     = []string{"ContractNumber", "Owner", "Status", "RegistrationDomain", "K8sClusterLimitTotal", "K8sClustersProvisioned"}
)

type ContractPrint struct {
	ContractNumber     int64  `json:"ContractNumber,omitempty"`
	Owner              string `json:"Owner,omitempty"`
	Status             string `json:"Status,omitempty"`
	RegistrationDomain string `json:"RegistrationDomain,omitempty"`
	// Contract Resource Limits
	CoresPerServer         int32 `json:"CoresPerServer,omitempty"`
	CoresPerContract       int32 `json:"CoresPerContract,omitempty"`
	CoresProvisioned       int32 `json:"CoresProvisioned,omitempty"`
	RamPerServer           int32 `json:"RamPerServer,omitempty"`
	RamPerContract         int32 `json:"RamPerContract,omitempty"`
	RamProvisioned         int32 `json:"RamProvisioned,omitempty"`
	HddLimitPerVolume      int64 `json:"HddLimitPerVolume,omitempty"`
	HddLimitPerContract    int64 `json:"HddLimitPerContract,omitempty"`
	HddVolumeProvisioned   int64 `json:"HddVolumeProvisioned,omitempty"`
	SsdLimitPerVolume      int64 `json:"SsdLimitPerVolume,omitempty"`
	SsdLimitPerContract    int64 `json:"SsdLimitPerContract,omitempty"`
	SsdVolumeProvisioned   int64 `json:"SsdVolumeProvisioned,omitempty"`
	ReservableIps          int32 `json:"ReservableIps,omitempty"`
	ReservedIpsOnContract  int32 `json:"ReservedIpsOnContract,omitempty"`
	ReservedIpsInUse       int32 `json:"ReservedIpsInUse,omitempty"`
	K8sClusterLimitTotal   int32 `json:"K8sClusterLimitTotal,omitempty"`
	K8sClustersProvisioned int32 `json:"K8sClustersProvisioned,omitempty"`
}

func getContractPrint(c *builder.CommandConfig, cs []resources.Contract, cols []string) printer.Result {
	r := printer.Result{}
	if c != nil {
		if cs != nil {
			r.OutputJSON = cs
			r.KeyValue = getContractsKVMaps(cs)
			r.Columns = cols
		}
	}
	return r
}

func getContractCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var contractCols []string
		columnsMap := map[string]string{
			"ContractNumber":         "ContractNumber",
			"Owner":                  "Owner",
			"Status":                 "Status",
			"RegistrationDomain":     "RegistrationDomain",
			"CoresPerServer":         "CoresPerServer",
			"CoresPerContract":       "CoresPerContract",
			"CoresProvisioned":       "CoresProvisioned",
			"RamPerServer":           "RamPerServer",
			"RamPerContract":         "RamPerContract",
			"RamProvisioned":         "RamProvisioned",
			"HddLimitPerVolume":      "HddLimitPerVolume",
			"HddLimitPerContract":    "HddLimitPerContract",
			"HddVolumeProvisioned":   "HddVolumeProvisioned",
			"SsdLimitPerVolume":      "SsdLimitPerVolume",
			"SsdLimitPerContract":    "SsdLimitPerContract",
			"SsdVolumeProvisioned":   "SsdVolumeProvisioned",
			"ReservableIps":          "ReservableIps",
			"ReservedIpsOnContract":  "ReservedIpsOnContract",
			"ReservedIpsInUse":       "ReservedIpsInUse",
			"K8sClusterLimitTotal":   "K8sClusterLimitTotal",
			"K8sClustersProvisioned": "K8sClustersProvisioned",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				contractCols = append(contractCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return contractCols
	} else {
		return defaultContractCols
	}
}

func getContract(c *resources.Contract) []resources.Contract {
	cs := make([]resources.Contract, 0)
	if c != nil {
		cs = append(cs, resources.Contract{Contract: c.Contract})
	}
	return cs
}

func getContractsKVMaps(cs []resources.Contract) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(cs))
	for _, c := range cs {
		o := getContractKVMap(c)
		out = append(out, o)
	}
	return out
}

func getContractKVMap(c resources.Contract) map[string]interface{} {
	var cPrint ContractPrint
	if properties, ok := c.GetPropertiesOk(); ok && properties != nil {
		if no, ok := properties.GetContractNumberOk(); ok && no != nil {
			cPrint.ContractNumber = *no
		}
		if owner, ok := properties.GetOwnerOk(); ok && owner != nil {
			cPrint.Owner = *owner
		}
		if status, ok := properties.GetStatusOk(); ok && status != nil {
			cPrint.Status = *status
		}
		if regDomain, ok := properties.GetRegDomainOk(); ok && regDomain != nil {
			cPrint.RegistrationDomain = *regDomain
		}
		if limits, ok := properties.GetResourceLimitsOk(); ok && limits != nil {
			cPrint = getResourceLimits(limits, cPrint)
		}
	}
	return structs.Map(cPrint)
}

func getResourceLimits(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	cPrint = getResourceLimitsCores(limits, cPrint)
	cPrint = getResourceLimitsRam(limits, cPrint)
	cPrint = getResourceLimitsHDD(limits, cPrint)
	cPrint = getResourceLimitsSSD(limits, cPrint)
	cPrint = getResourceLimitsIPS(limits, cPrint)
	cPrint = getResourceLimitsK8S(limits, cPrint)
	return cPrint
}

func getResourceLimitsCores(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if coresServer, ok := limits.GetCoresPerServerOk(); ok && coresServer != nil {
		cPrint.CoresPerServer = *coresServer
	}
	if coresContract, ok := limits.GetCoresPerContractOk(); ok && coresContract != nil {
		cPrint.CoresPerContract = *coresContract
	}
	if coresProvisioned, ok := limits.GetCoresProvisionedOk(); ok && coresProvisioned != nil {
		cPrint.CoresProvisioned = *coresProvisioned
	}
	return cPrint
}

func getResourceLimitsRam(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if ramServer, ok := limits.GetRamPerServerOk(); ok && ramServer != nil {
		cPrint.RamPerServer = *ramServer
	}
	if ramContract, ok := limits.GetRamPerContractOk(); ok && ramContract != nil {
		cPrint.RamPerContract = *ramContract
	}
	if ramProvisioned, ok := limits.GetRamProvisionedOk(); ok && ramProvisioned != nil {
		cPrint.RamProvisioned = *ramProvisioned
	}
	return cPrint
}

func getResourceLimitsHDD(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if hddVolume, ok := limits.GetHddLimitPerVolumeOk(); ok && hddVolume != nil {
		cPrint.HddLimitPerVolume = *hddVolume
	}
	if hddVolumeContract, ok := limits.GetHddLimitPerContractOk(); ok && hddVolumeContract != nil {
		cPrint.HddLimitPerContract = *hddVolumeContract
	}
	if hddVolumeProvisioned, ok := limits.GetHddVolumeProvisionedOk(); ok && hddVolumeProvisioned != nil {
		cPrint.HddVolumeProvisioned = *hddVolumeProvisioned
	}
	return cPrint
}

func getResourceLimitsSSD(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if ssdVolume, ok := limits.GetSsdLimitPerVolumeOk(); ok && ssdVolume != nil {
		cPrint.SsdLimitPerVolume = *ssdVolume
	}
	if ssdVolumeContract, ok := limits.GetSsdLimitPerContractOk(); ok && ssdVolumeContract != nil {
		cPrint.SsdLimitPerContract = *ssdVolumeContract
	}
	if ssdVolumeProvisioned, ok := limits.GetSsdVolumeProvisionedOk(); ok && ssdVolumeProvisioned != nil {
		cPrint.SsdVolumeProvisioned = *ssdVolumeProvisioned
	}
	return cPrint
}

func getResourceLimitsIPS(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if reservableIps, ok := limits.GetReservableIpsOk(); ok && reservableIps != nil {
		cPrint.ReservableIps = *reservableIps
	}
	if reservedIpsContract, ok := limits.GetReservedIpsOnContractOk(); ok && reservedIpsContract != nil {
		cPrint.ReservedIpsOnContract = *reservedIpsContract
	}
	if reservedIpsUse, ok := limits.GetReservedIpsInUseOk(); ok && reservedIpsUse != nil {
		cPrint.ReservedIpsInUse = *reservedIpsUse
	}
	return cPrint
}

func getResourceLimitsK8S(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if clusterTotal, ok := limits.GetK8sClusterLimitTotalOk(); ok && clusterTotal != nil {
		cPrint.K8sClusterLimitTotal = *clusterTotal
	}
	if clusterProvisioned, ok := limits.GetK8sClustersProvisionedOk(); ok && clusterProvisioned != nil {
		cPrint.K8sClustersProvisioned = *clusterProvisioned
	}
	return cPrint
}
