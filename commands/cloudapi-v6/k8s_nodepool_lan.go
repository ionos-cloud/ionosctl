package commands

import (
	"context"
	"errors"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sNodePoolLanCmd() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Short:            "Kubernetes NodePool LAN Operations",
			Long:             "The sub-commands of `ionosctl k8s nodepool lan` allow you to list, add, remove Kubernetes Node Pool LANs.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "nodepool",
		Resource:   "lan",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes NodePool LANs",
		LongDesc:   "Use this command to get a list of all contained NodePool LANs in a selected Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id",
		Example:    listK8sNodePoolLanExample,
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolLanList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgK8sNodePoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.ArgCols, "", defaultK8sNodePoolLanCols, printer.ColsMessage(defaultK8sNodePoolLanCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodePoolLanCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	/*
		Add Command
	*/
	add := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "nodepool",
		Resource:  "lan",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Kubernetes NodePool LAN",
		LongDesc: `Use this command to add a Node Pool LAN into an existing Node Pool.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run a command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id`,
		Example:    addK8sNodePoolLanExample,
		PreCmdRun:  PreRunK8sClusterNodePoolLanIds,
		CmdRun:     RunK8sNodePoolLanAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgK8sNodePoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "The unique LAN Id of existing LANs to be attached to worker Nodes", core.RequiredFlagOption())
	add.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	add.AddStringFlag(cloudapiv6.ArgNetwork, "", "", "IPv4 or IPv6 CIDR to be routed via the interface. Must be set with --gateway-ip flag")
	add.AddIpFlag(cloudapiv6.ArgGatewayIp, "", nil, "IPv4 or IPv6 Gateway IP for the route. Must be set with --network flag")
	add.AddStringSliceFlag(cloudapiv6.ArgCols, "", defaultK8sNodePoolLanCols, printer.ColsMessage(defaultK8sNodePoolLanCols))
	_ = add.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodePoolLanCols, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "nodepool",
		Resource:  "lan",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a Kubernetes NodePool LAN",
		LongDesc: `This command removes a Kubernetes Node Pool LAN from a Node Pool.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id`,
		Example:    removeK8sNodePoolLanExample,
		PreCmdRun:  PreRunK8sClusterNodePoolLanRemove,
		CmdRun:     RunK8sNodePoolLanRemove,
		InitClient: true,
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgK8sNodePoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "The unique LAN Id of existing LANs to be detached from worker Nodes", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all FK8s Nodepool Lans.")
	removeCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return k8sCmd
}

func PreRunK8sClusterNodePoolLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgK8sClusterId, cloudapiv6.ArgK8sNodePoolId, cloudapiv6.ArgLanId)
}

func PreRunK8sClusterNodePoolLanRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgK8sClusterId, cloudapiv6.ArgK8sNodePoolId, cloudapiv6.ArgLanId},
		[]string{cloudapiv6.ArgK8sClusterId, cloudapiv6.ArgK8sNodePoolId, cloudapiv6.ArgAll},
	)
}

func RunK8sNodePoolLanList(c *core.CommandConfig) error {
	k8ss, resp, err := c.CloudApiV6Services.K8s().GetNodePool(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId)),
		resources.QueryParams{},
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if properties, ok := k8ss.GetPropertiesOk(); ok && properties != nil {
		if lans, ok := properties.GetLansOk(); ok && lans != nil {
			return c.Printer.Print(getK8sNodePoolLanPrint(c, getK8sNodePoolLans(lans)))
		} else {
			return errors.New("error getting node pool lans")
		}
	} else {
		return errors.New("error getting node pool properties")
	}
}

func RunK8sNodePoolLanAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	ng, _, err := c.CloudApiV6Services.K8s().GetNodePool(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId)),
		queryParams,
	)
	if err != nil {
		return err
	}
	input := getNewK8sNodePoolLanInfo(c, ng)
	ngNew, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolLanPrint(c, getK8sNodePoolLansForPut(ngNew)))
}

func RunK8sNodePoolLanRemove(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllK8sNodePoolsLans(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		clusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId))
		nodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId))
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove node pool lan"); err != nil {
			return err
		}
		ng, _, err := c.CloudApiV6Services.K8s().GetNodePool(clusterId, nodePoolId, queryParams)
		if err != nil {
			return err
		}
		input := removeK8sNodePoolLanInfo(c, ng)
		_, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(clusterId, nodePoolId, input, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		return c.Printer.Warn("Status: Command node pool lan remove has been successfully executed")
	}
}

func RemoveAllK8sNodePoolsLans(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	clusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId))
	nodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId))
	c.Printer.Verbose("K8sCluster ID: %v", clusterId)
	c.Printer.Verbose("K8sNodePool ID: %v", nodePoolId)
	c.Printer.Verbose("Getting K8sNodePool Lans...")
	k8sNodepool, resp, err := c.CloudApiV6Services.K8s().GetNodePool(clusterId, nodePoolId, cloudapiv6.ParentResourceQueryParams)
	if err != nil {
		return err
	}
	if nodePoolProperties, ok := k8sNodepool.GetPropertiesOk(); ok && nodePoolProperties != nil {
		if lans, ok := nodePoolProperties.GetLansOk(); ok && lans != nil {
			if len(*lans) > 0 {
				_ = c.Printer.Warn("K8s NodePool Lans to be removed:")
				for _, lan := range *lans {
					if id, ok := lan.GetIdOk(); ok && id != nil {
						_ = c.Printer.Warn("K8s NodePool Lan Id: " + string(*id))
					}
				}
			} else {
				return errors.New("no Lans found")
			}
		} else {
			return errors.New("could not get Lans items")
		}
	}
	if err = utils.AskForConfirm(c.Stdin, c.Printer, "remove all the K8sNodePool Lans"); err != nil {
		return err
	}
	c.Printer.Verbose("Removing all the K8sNodePool Lans...")
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if properties, ok := k8sNodepool.GetPropertiesOk(); ok && properties != nil {
		if n, ok := properties.GetNodeCountOk(); ok && n != nil {
			propertiesUpdated.SetNodeCount(*n)
		}
		if n, ok := properties.GetAutoScalingOk(); ok && n != nil {
			propertiesUpdated.SetAutoScaling(*n)
		}
		if n, ok := properties.GetMaintenanceWindowOk(); ok && n != nil {
			propertiesUpdated.SetMaintenanceWindow(*n)
		}
		if n, ok := properties.GetK8sVersionOk(); ok && n != nil {
			propertiesUpdated.SetK8sVersion(*n)
		}
		newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
		propertiesUpdated.SetLans(newLans)
		k8sNodePoolUpdated := resources.K8sNodePoolForPut{
			KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
				Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
			},
		}
		_, resp, err = c.CloudApiV6Services.K8s().UpdateNodePool(clusterId, nodePoolId, k8sNodePoolUpdated, queryParams)
		if resp != nil {
			c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func getNewK8sNodePoolLanInfo(c *core.CommandConfig, oldNg *resources.K8sNodePool) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
		if n, ok := properties.GetNodeCountOk(); ok && n != nil {
			propertiesUpdated.SetNodeCount(*n)
		}
		if n, ok := properties.GetAutoScalingOk(); ok && n != nil {
			propertiesUpdated.SetAutoScaling(*n)
		}
		if n, ok := properties.GetMaintenanceWindowOk(); ok && n != nil {
			propertiesUpdated.SetMaintenanceWindow(*n)
		}
		if n, ok := properties.GetK8sVersionOk(); ok && n != nil {
			propertiesUpdated.SetK8sVersion(*n)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			// Append existing LANs
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}
			// Add new LANs
			lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
			dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
			newLan := ionoscloud.KubernetesNodePoolLan{
				Id:   &lanId,
				Dhcp: &dhcp,
			}
			c.Printer.Verbose("Adding a Kubernetes NodePool LAN with id: %v and dhcp: %v", lanId, dhcp)
			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork)) &&
				viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgGatewayIp)) {
				newRoute := ionoscloud.KubernetesNodePoolLanRoutes{}
				if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork)) {
					network := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork))
					newRoute.SetNetwork(network)
					c.Printer.Verbose("Property Network set: %v", network)
				}
				if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgGatewayIp)) {
					gatewayIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGatewayIp))
					newRoute.SetGatewayIp(gatewayIp)
					c.Printer.Verbose("Property GatewayIp set: %v", gatewayIp)
				}
				newLan.SetRoutes([]ionoscloud.KubernetesNodePoolLanRoutes{newRoute})
			}
			newLans = append(newLans, newLan)
			propertiesUpdated.SetLans(newLans)
		}
	}
	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

func removeK8sNodePoolLanInfo(c *core.CommandConfig, oldNg *resources.K8sNodePool) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
		if n, ok := properties.GetNodeCountOk(); ok && n != nil {
			propertiesUpdated.SetNodeCount(*n)
		}
		if n, ok := properties.GetAutoScalingOk(); ok && n != nil {
			propertiesUpdated.SetAutoScaling(*n)
		}
		if n, ok := properties.GetMaintenanceWindowOk(); ok && n != nil {
			propertiesUpdated.SetMaintenanceWindow(*n)
		}
		if n, ok := properties.GetK8sVersionOk(); ok && n != nil {
			propertiesUpdated.SetK8sVersion(*n)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
			lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			c.Printer.Verbose("Removing a Kubernetes NodePool LAN with id: %v", lanId)
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					if id, ok := existingLan.GetIdOk(); ok && id != nil {
						if *id != lanId {
							newLans = append(newLans, existingLan)
						}
					}
				}
			}
			propertiesUpdated.SetLans(newLans)
		}
	}
	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

// Output Printing

var defaultK8sNodePoolLanCols = []string{"LanId", "Dhcp", "RoutesNetwork", "RoutesGatewayIp"}

type K8sNodePoolLanPrint struct {
	LanId           int32    `json:"LanId,omitempty"`
	Dhcp            bool     `json:"Dhcp,omitempty"`
	RoutesNetwork   []string `json:"RoutesNetwork,omitempty"`
	RoutesGatewayIp []string `json:"RoutesGatewayIp,omitempty"`
}

func getK8sNodePoolLanPrint(c *core.CommandConfig, k8ss []resources.K8sNodePoolLan) printer.Result {
	r := printer.Result{}
	if c != nil {
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sNodePoolLansKVMaps(k8ss)
			r.Columns = printer.GetHeadersAllDefault(defaultK8sNodePoolLanCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)))
		}
	}
	return r
}

func getK8sNodePoolLansForPut(ng *resources.K8sNodePool) []resources.K8sNodePoolLan {
	ss := make([]resources.K8sNodePoolLan, 0)
	if ng != nil {
		if properties, ok := ng.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					ss = append(ss, resources.K8sNodePoolLan{
						KubernetesNodePoolLan: lanItem,
					})
				}
			}
		}
	}
	return ss
}

func getK8sNodePoolLans(k8ss *[]ionoscloud.KubernetesNodePoolLan) []resources.K8sNodePoolLan {
	u := make([]resources.K8sNodePoolLan, 0)
	if k8ss != nil {
		for _, item := range *k8ss {
			u = append(u, resources.K8sNodePoolLan{KubernetesNodePoolLan: item})
		}
	}
	return u
}

func getK8sNodePoolLansKVMaps(us []resources.K8sNodePoolLan) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint K8sNodePoolLanPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.LanId = *id
		}
		if dhcp, ok := u.GetDhcpOk(); ok && dhcp != nil {
			uPrint.Dhcp = *dhcp
		}
		if routes, ok := u.GetRoutesOk(); ok && routes != nil {
			newRoutesNetwork := make([]string, 0)
			newRoutesGatewayIp := make([]string, 0)
			for _, route := range *routes {
				if net, ok := route.GetNetworkOk(); ok && net != nil {
					newRoutesNetwork = append(newRoutesNetwork, *net)
				}
				if ip, ok := route.GetGatewayIpOk(); ok && ip != nil {
					newRoutesGatewayIp = append(newRoutesGatewayIp, *ip)
				}
			}
			uPrint.RoutesNetwork = newRoutesNetwork
			uPrint.RoutesGatewayIp = newRoutesGatewayIp
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}
