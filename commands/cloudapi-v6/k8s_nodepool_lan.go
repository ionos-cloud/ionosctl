package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allK8sNodePoolLanJSONPaths = map[string]string{
		"LanId":           "id",
		"Dhcp":            "dhcp",
		"RoutesNetwork":   "routes.*.network",
		"RoutesGatewayIp": "routes.*.gatewayIp",
	}

	defaultK8sNodePoolLanCols = []string{"LanId", "Dhcp", "RoutesNetwork", "RoutesGatewayIp"}
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
	list.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(constants.FlagNodepoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.ArgCols, "", defaultK8sNodePoolLanCols, printer.ColsMessage(defaultK8sNodePoolLanCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodePoolLanCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
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
	add.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(constants.FlagNodepoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "The unique LAN Id of existing LANs to be attached to worker Nodes", core.RequiredFlagOption())
	add.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	add.AddStringSliceFlag(cloudapiv6.ArgNetwork, "", nil, "Slice of IPv4 or IPv6 CIDRs to be routed via the interface. Must contain same number of arguments as --gateway-ip flag")
	add.AddStringSliceFlag(cloudapiv6.ArgGatewayIp, "", nil, "Slice of IPv4 or IPv6 Gateway IPs for the routes. Must contain same number of arguments as --network flag")
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
	removeCmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(constants.FlagNodepoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "The unique LAN Id of existing LANs to be detached from worker Nodes", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all FK8s Nodepool Lans.")
	removeCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return k8sCmd
}

func PreRunK8sClusterNodePoolLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgLanId)
}

func PreRunK8sClusterNodePoolLanRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgLanId},
		[]string{constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgAll},
	)
}

func RunK8sNodePoolLanList(c *core.CommandConfig) error {
	k8ss, resp, err := c.CloudApiV6Services.K8s().GetNodePool(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		resources.QueryParams{},
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := k8ss.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting node pool properties")
	}

	lans, ok := properties.GetLansOk()
	if !ok || lans == nil {
		return fmt.Errorf("error getting node pool lans")
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allK8sNodePoolLanJSONPaths, lans,
		printer.GetHeadersAllDefault(defaultK8sNodePoolLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunK8sNodePoolLanAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	ng, _, err := c.CloudApiV6Services.K8s().GetNodePool(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		queryParams,
	)
	if err != nil {
		return err
	}

	input := getNewK8sNodePoolLanInfo(c, ng)
	ngNew, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allK8sNodePoolLanJSONPaths, getK8sNodePoolLansForPut(ngNew),
		printer.GetHeadersAllDefault(defaultK8sNodePoolLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
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

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Kubernetes Node Pool Lans successfully deleted"))
		return nil
	}

	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	nodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	if !confirm.Ask("remove node pool lan") {
		return nil
	}

	ng, _, err := c.CloudApiV6Services.K8s().GetNodePool(clusterId, nodePoolId, queryParams)
	if err != nil {
		return err
	}

	input := removeK8sNodePoolLanInfo(c, ng)
	_, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(clusterId, nodePoolId, input, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Kubernetes Node Pool Lan successfully deleted"))
	return nil
}

func RemoveAllK8sNodePoolsLans(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	nodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("K8sCluster ID: %v", clusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("K8sNodePool ID: %v", nodePoolId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting K8sNodePool Lans..."))

	k8sNodepool, resp, err := c.CloudApiV6Services.K8s().GetNodePool(clusterId, nodePoolId, cloudapiv6.ParentResourceQueryParams)
	if err != nil {
		return err
	}

	nodePoolProperties, ok := k8sNodepool.GetPropertiesOk()
	if !ok || nodePoolProperties == nil {
		return fmt.Errorf("could not get Node Pool properties")
	}

	lans, ok := nodePoolProperties.GetLansOk()
	if !ok || lans == nil {
		return fmt.Errorf("could not get Lans items")
	}

	if len(*lans) <= 0 {
		return fmt.Errorf("no Lans found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("K8s NodePool Lans to be removed:"))
	for _, lan := range *lans {
		if id, ok := lan.GetIdOk(); ok && id != nil {
			fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("K8s NodePool Lan Id: "+string(*id)))
		}
	}

	if !confirm.Ask("remove all the K8sNodePool Lans") {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing all the K8sNodePool Lans..."))

	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if n, ok := nodePoolProperties.GetNodeCountOk(); ok && n != nil {
		propertiesUpdated.SetNodeCount(*n)
	}

	if n, ok := nodePoolProperties.GetAutoScalingOk(); ok && n != nil {
		propertiesUpdated.SetAutoScaling(*n)
	}

	if n, ok := nodePoolProperties.GetMaintenanceWindowOk(); ok && n != nil {
		propertiesUpdated.SetMaintenanceWindow(*n)
	}

	if n, ok := nodePoolProperties.GetK8sVersionOk(); ok && n != nil {
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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
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

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Adding a Kubernetes NodePool LAN with id: %v and dhcp: %v", lanId, dhcp))

			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork)) {
				network := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork))
				gatewayIp := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgGatewayIp))

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Network set: %v", network))
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property GatewayIp set: %v", gatewayIp))

				if len(network) != len(gatewayIp) {
					die.Die(fmt.Sprintf("Flags %s, %s have different number of arguments, must be the same", cloudapiv6.ArgNetwork, cloudapiv6.ArgGatewayIp))
				}

				routes := make([]ionoscloud.KubernetesNodePoolLanRoutes, 0)
				for i, net := range network {
					routes = append(routes,
						ionoscloud.KubernetesNodePoolLanRoutes{
							Network:   pointer.From(net), // Copy the loop variable and take its address. See #289 - always same address would be used
							GatewayIp: &gatewayIp[i],
						},
					)
				}

				newLan.SetRoutes(routes)
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

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing a Kubernetes NodePool LAN with id: %v", lanId))

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

func getK8sNodePoolLansForPut(ng *resources.K8sNodePool) []ionoscloud.KubernetesNodePoolLan {
	ss := make([]ionoscloud.KubernetesNodePoolLan, 0)

	if ng != nil {
		if properties, ok := ng.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					ss = append(ss, lanItem)
				}
			}
		}
	}

	return ss
}
