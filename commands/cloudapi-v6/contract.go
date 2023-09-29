package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allContractJSONPaths = map[string]string{
		"ContractNumber":         "properties.contractNumber",
		"Owner":                  "properties.owner",
		"Status":                 "properties.status",
		"RegistrationDomain":     "properties.regDomain",
		"CoresPerServer":         "properties.resourceLimits.coresPerServer",
		"CoresPerContract":       "properties.resourceLimits.coresPerContract",
		"CoresProvisioned":       "properties.resourceLimits.coresProvisioned",
		"RamPerServer":           "properties.resourceLimits.ramPerServer",
		"RamPerContract":         "properties.resourceLimits.ramPerContract",
		"RamProvisioned":         "properties.resourceLimits.ramProvisioned",
		"HddLimitPerVolume":      "properties.resourceLimits.hddLimitPerVolume",
		"HddLimitPerContract":    "properties.resourceLimits.hddLimitPerContract",
		"HddVolumeProvisioned":   "properties.resourceLimits.hddVolumeProvisioned",
		"SsdLimitPerVolume":      "properties.resourceLimits.ssdLimitPerVolume",
		"SsdLimitPerContract":    "properties.resourceLimits.ssdLimitPerContract",
		"SsdVolumeProvisioned":   "properties.resourceLimits.ssdVolumeProvisioned",
		"DasVolumeProvisioned":   "properties.resourceLimits.dasVolumeProvisioned",
		"ReservableIps":          "properties.resourceLimits.reservableIps",
		"ReservedIpsOnContract":  "properties.resourceLimits.reservedIpsOnContract",
		"ReservedIpsInUse":       "properties.resourceLimits.reserverIpsInUse",
		"K8sClusterLimitTotal":   "k8sClusterLimitTotal",
		"K8sClustersProvisioned": "k8sClustersProvisioned",
		"NlbLimitTotal":          "properties.resourceLimits.nlbLimitTotal",
		"NlbProvisioned":         "properties.resourceLimits.nlbProvisioned",
		"NatGatewayLimitTotal":   "properties.resourceLimits.natGatewayLimitTotal",
		"NatGatewayProvisioned":  "properties.resourceLimits.natGatewayProvisioned",
	}

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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultContractCols, tabheaders.ColsMessage(allContractCols))
	_ = viper.BindPFlag(core.GetFlagName(contractCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = contractCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	return contractCmd
}

func RunContractGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Contract with resource limits: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits))))

	contractResource, resp, err := c.CloudApiV6Services.Contracts().Get(queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	var out string

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)) {
		switch strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits))) {
		case "CORES":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractCoresCols)
		case "RAM":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractRamCols)
		case "HDD":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractHddCols)
		case "SSD":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractSsdCols)
		case "DAS":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractDasCols)
		case "IPS":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractIpsCols)
		case "K8S":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractK8sCols)
		case "NLB":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractNlbCols)
		case "NAT":
			out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts, contractNatCols)
		}
	} else {
		out, err = jsontabwriter.GenerateOutput("items", allContractJSONPaths, contractResource.Contracts,
			tabheaders.GetHeaders(allContractCols, defaultContractCols, cols))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
