package commands

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ContractCmd() *core.Command {
	ctx := context.TODO()
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultContractCols, printer.ColsMessage(allContractCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(contractCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = contractCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allContractCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, contractCmd, core.CommandBuilder{
		Namespace:  "contract",
		Resource:   "contract",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get information about the Contract Resources on your account",
		LongDesc:   "Use this command to get information about the Contract Resources on your account. Use `--resource-limits` flag to see specific Contract Resources Limits.",
		Example:    getContractExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunContractGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.ArgResourceLimits, "", "", "Specify Resource Limits to see details about it")
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceLimits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CORES", "RAM", "HDD", "SSD", "DAS", "IPS", "K8S", "NLB", "NAT"}, cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return contractCmd
}

func RunContractGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Contract with resource limits: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)))
	contractResource, resp, err := c.CloudApiV6Services.Contracts().Get()
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)) {
		switch strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits))) {
		case "CORES":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractCoresCols))
		case "RAM":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractRamCols))
		case "HDD":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractHddCols))
		case "SSD":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractSsdCols))
		case "DAS":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractDasCols))
		case "IPS":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractIpsCols))
		case "K8S":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractK8sCols))
		case "NLB":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractNlbCols))
		case "NAT":
			return c.Printer.Print(getContractPrint(c, getContract(contractResource), contractNatCols))
		}
	}
	return c.Printer.Print(getContractPrint(c, getContract(contractResource), getContractCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())))
}

// Output Printing

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
	DasVolumeProvisioned   int64 `json:"DasVolumeProvisioned,omitempty"`
	ReservableIps          int32 `json:"ReservableIps,omitempty"`
	ReservedIpsOnContract  int32 `json:"ReservedIpsOnContract,omitempty"`
	ReservedIpsInUse       int32 `json:"ReservedIpsInUse,omitempty"`
	K8sClusterLimitTotal   int32 `json:"K8sClusterLimitTotal,omitempty"`
	K8sClustersProvisioned int32 `json:"K8sClustersProvisioned,omitempty"`
	NlbLimitTotal          int32 `json:"NlbLimitTotal,omitempty"`
	NlbProvisioned         int32 `json:"NlbProvisioned,omitempty"`
	NatGatewayLimitTotal   int32 `json:"NatGatewayLimitTotal,omitempty"`
	NatGatewayProvisioned  int32 `json:"NatGatewayProvisioned,omitempty"`
}

func getContractPrint(c *core.CommandConfig, cs []resources.Contract, cols []string) printer.Result {
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
			"DasVolumeProvisioned":   "DasVolumeProvisioned",
			"ReservableIps":          "ReservableIps",
			"ReservedIpsOnContract":  "ReservedIpsOnContract",
			"ReservedIpsInUse":       "ReservedIpsInUse",
			"K8sClusterLimitTotal":   "K8sClusterLimitTotal",
			"K8sClustersProvisioned": "K8sClustersProvisioned",
			"NlbLimitTotal":          "NlbLimitTotal",
			"NlbProvisioned":         "NlbProvisioned",
			"NatGatewayLimitTotal":   "NatGatewayLimitTotal",
			"NatGatewayProvisioned":  "NatGatewayProvisioned",
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

func getContract(c resources.Contracts) []resources.Contract {
	cs := make([]resources.Contract, 0)
	if items, ok := c.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			cs = append(cs, resources.Contract{Contract: item})
		}
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
	cPrint = getResourceLimitsDAS(limits, cPrint)
	cPrint = getResourceLimitsIPS(limits, cPrint)
	cPrint = getResourceLimitsK8S(limits, cPrint)
	cPrint = getResourceLimitsNatGateway(limits, cPrint)
	cPrint = getResourceLimitsNlb(limits, cPrint)
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

func getResourceLimitsDAS(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if dasVolume, ok := limits.GetDasVolumeProvisionedOk(); ok && dasVolume != nil {
		cPrint.DasVolumeProvisioned = *dasVolume
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

func getResourceLimitsNatGateway(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if natTotal, ok := limits.GetNatGatewayLimitTotalOk(); ok && natTotal != nil {
		cPrint.NatGatewayLimitTotal = *natTotal
	}
	if natProvisioned, ok := limits.GetNatGatewayProvisionedOk(); ok && natProvisioned != nil {
		cPrint.NatGatewayProvisioned = *natProvisioned
	}
	return cPrint
}

func getResourceLimitsNlb(limits *ionoscloud.ResourceLimits, cPrint ContractPrint) ContractPrint {
	if nlbTotal, ok := limits.GetNlbLimitTotalOk(); ok && nlbTotal != nil {
		cPrint.NlbLimitTotal = *nlbTotal
	}
	if nlbProvisioned, ok := limits.GetNlbProvisionedOk(); ok && nlbProvisioned != nil {
		cPrint.NlbProvisioned = *nlbProvisioned
	}
	return cPrint
}
