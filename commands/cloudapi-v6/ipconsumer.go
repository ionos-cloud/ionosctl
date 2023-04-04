package commands

import (
	"context"
	"errors"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultIpConsumerCols, printer.ColsMessage(allIpConsumerCols))
	_ = viper.BindPFlag(core.GetFlagName(resourceCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	listResources.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = listResources.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listResources.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	listResources.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	listResources.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	return resourceCmd
}

func RunIpConsumersList(c *core.CommandConfig) error {
	ipBlock, resp, err := c.CloudApiV6Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)), resources.QueryParams{})
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if properties, ok := ipBlock.GetPropertiesOk(); ok && properties != nil {
		if ipCons, ok := properties.GetIpConsumersOk(); ok && ipCons != nil {
			ipsConsumers := make([]resources.IpConsumer, 0)
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
			r.Columns = printer.GetHeaders(allIpConsumerCols, defaultIpConsumerCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)))
		}
	}
	return r
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
