package cloudapi_v5

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IpconsumerCmd() *core.Command {
	ctx := context.TODO()
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipconsumer",
			Aliases:          []string{"ipc"},
			Short:            "Ip Consumer Operations",
			Long:             "The sub-command of `ionosctl ipconsumer` allows you to list information about where each IP address from an IpBlock is being used.",
			TraverseChildren: true,
		},
	}
	globalFlags := resourceCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultIpConsumerCols, printer.ColsMessage(allIpConsumerCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(resourceCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allIpConsumerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	listResources := core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "ipconsumer",
		Resource:   "ipconsumer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IpConsumers",
		LongDesc:   "Use this command to get a list of Resources where each IP address from an IpBlock is being used.\n\nRequired values to run command:\n\n* IpBlock Id",
		Example:    listIpConsumersExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpConsumersList,
		InitClient: true,
	})
	listResources.AddStringFlag(cloudapiv5.ArgIpBlockId, "", "", cloudapiv5.IpBlockId, core.RequiredFlagOption())
	_ = listResources.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listResources.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return resourceCmd
}

func RunIpConsumersList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting IpBlock with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId)))
	ipBlock, resp, err := c.CloudApiV5Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if properties, ok := ipBlock.GetPropertiesOk(); ok && properties != nil {
		if ipCons, ok := properties.GetIpConsumersOk(); ok && ipCons != nil {
			ipsConsumers := make([]resources.IpConsumer, 0)
			c.Printer.Verbose("Getting IpConsumers from IpBlock...")
			for _, ip := range *ipCons {
				ipsConsumers = append(ipsConsumers, resources.IpConsumer{IpConsumer: ip})
			}
			return c.Printer.Print(getIpConsumerPrint(c, ipsConsumers))
		} else {
			return errors.New("error getting ipconsumers")
		}
	} else {
		return errors.New("error getting ipblock properties")
	}
}

// Output Printing

var (
	defaultIpConsumerCols = []string{"Ip", "NicId", "ServerId", "DatacenterId", "K8sNodePoolId", "K8sClusterId"}
	allIpConsumerCols     = []string{"Ip", "Mac", "NicId", "ServerId", "ServerName", "DatacenterId", "DatacenterName", "K8sNodePoolId", "K8sClusterId"}
)

type IpConsumerPrint struct {
	Ip             string `json:"Ip,omitempty"`
	Mac            string `json:"Mac,omitempty"`
	NicId          string `json:"NicId,omitempty"`
	ServerId       string `json:"ServerId,omitempty"`
	ServerName     string `json:"ServerName,omitempty"`
	DatacenterId   string `json:"DatacenterId,omitempty"`
	DatacenterName string `json:"DatacenterName,omitempty"`
	K8sNodePoolId  string `json:"K8sNodePoolId,omitempty"`
	K8sClusterId   string `json:"K8sClusterId,omitempty"`
}

func getIpConsumerPrint(c *core.CommandConfig, groups []resources.IpConsumer) printer.Result {
	r := printer.Result{}
	if c != nil {
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getIpConsumersKVMaps(groups)
			r.Columns = getIpConsumerCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getIpConsumerCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var groupCols []string
		columnsMap := map[string]string{
			"Ip":             "Ip",
			"Mac":            "Mac",
			"NicId":          "NicId",
			"ServerId":       "ServerId",
			"ServerName":     "ServerName",
			"DatacenterId":   "DatacenterId",
			"DatacenterName": "DatacenterName",
			"K8sNodePoolId":  "K8sNodePoolId",
			"K8sClusterId":   "K8sClusterId",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				groupCols = append(groupCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return groupCols
	} else {
		return defaultIpConsumerCols
	}
}

func getIpConsumersKVMaps(rs []resources.IpConsumer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(rs))
	for _, r := range rs {
		var rPrint IpConsumerPrint
		if ip, ok := r.GetIpOk(); ok && ip != nil {
			rPrint.Ip = *ip
		}
		if mac, ok := r.GetMacOk(); ok && mac != nil {
			rPrint.Mac = *mac
		}
		if nicId, ok := r.GetNicIdOk(); ok && nicId != nil {
			rPrint.NicId = *nicId
		}
		if serverId, ok := r.GetServerIdOk(); ok && serverId != nil {
			rPrint.ServerId = *serverId
		}
		if serverName, ok := r.GetServerNameOk(); ok && serverName != nil {
			rPrint.ServerName = *serverName
		}
		if datacenterId, ok := r.GetDatacenterIdOk(); ok && datacenterId != nil {
			rPrint.DatacenterId = *datacenterId
		}
		if datacenterName, ok := r.GetDatacenterNameOk(); ok && datacenterName != nil {
			rPrint.DatacenterName = *datacenterName
		}
		if nodepoolId, ok := r.GetK8sNodePoolUuidOk(); ok && nodepoolId != nil {
			rPrint.K8sNodePoolId = *nodepoolId
		}
		if clusterId, ok := r.GetK8sClusterUuidOk(); ok && clusterId != nil {
			rPrint.K8sClusterId = *clusterId
		}
		o := structs.Map(rPrint)
		out = append(out, o)
	}
	return out
}
