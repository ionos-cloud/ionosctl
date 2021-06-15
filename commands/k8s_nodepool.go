package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8sNodePool() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nodepool",
			Aliases:          []string{"np"},
			Short:            "Kubernetes NodePool Operations",
			Long:             "The sub-commands of `ionosctl k8s nodepool` allow you to list, get, create, update, delete Kubernetes NodePools.",
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultK8sNodePoolCols, utils.ColsMessage(allK8sNodePoolCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(k8sCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sNodePoolCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "nodepool",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes NodePools",
		LongDesc:   "Use this command to get a list of all contained NodePools in a selected Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    listK8sNodePoolsExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sNodePoolList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "nodepool",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Kubernetes NodePool",
		LongDesc:   "Use this command to retrieve details about a specific NodePool from an existing Kubernetes Cluster. You can wait for the Node Pool to be in \"ACTIVE\" state using `--wait-for-state` flag together with `--timeout` option.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id",
		Example:    getK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgK8sNodePoolId, config.ArgIdShort, "", config.RequiredFlagK8sNodePoolId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified NodePool to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "nodepool",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Kubernetes NodePool",
		LongDesc: `Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center. Regarding the name for the Kubernetes NodePool, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. Regarding the Kubernetes Version for the NodePool, if not set via flag, it will be used the default one: ` + "`" + `ionosctl k8s version get` + "`" + `.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run a command:

* K8s Cluster Id
* Datacenter Id
* Name`,
		Example:    createK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterDcIdsNodePoolName,
		CmdRun:     RunK8sNodePoolCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name for the K8s NodePool "+config.RequiredFlag)
	create.AddStringFlag(config.ArgK8sVersion, "", "", "The K8s version for the NodePool. If not set, it will be used the default one")
	create.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability")
	create.AddIntFlag(config.ArgCores, "", 2, "The total number of cores for the Node")
	create.AddStringFlag(config.ArgRam, "", strconv.Itoa(2048), "RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgCpuFamily, "", config.DefaultServerCPUFamily, "CPU Type")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgAvailabilityZone, config.ArgAvailabilityZoneShort, "AUTO", "The compute Availability Zone in which the Node should exist")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgStorageType, "", "HDD", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgStorageSize, "", 10, "The total allocated storage capacity of a Node")
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for NodePool/Request [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "nodepool",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Kubernetes NodePool",
		LongDesc: `Use this command to update the number of worker Nodes, the minimum and maximum number of worker Nodes, the add labels, annotations, to update the maintenance day and time, to attach private LANs to a Node Pool within an existing Kubernetes Cluster. You can also add reserved public IP addresses to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. The array must contain one extra IP than maximum number of Nodes could be. (nodeCount+1 if fixed node amount or maxNodeCount+1 if auto scaling is used) The extra provided IP Will be used during rebuilding of Nodes.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`,
		Example:    updateK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgK8sVersion, "", "", "The K8s version for the NodePool. K8s version downgrade is not supported")
	update.AddIntFlag(config.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the NodePool should contain")
	update.AddIntFlag(config.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count")
	update.AddIntFlag(config.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count")
	update.AddStringFlag(config.ArgLabelKey, "", "", "Label key. Must be set together with --label-value")
	update.AddStringFlag(config.ArgLabelValue, "", "", "Label value. Must be set together with --label-key")
	update.AddStringFlag(config.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value")
	update.AddStringFlag(config.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key")
	update.AddStringFlag(config.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringSliceFlag(config.ArgPublicIps, "", []string{""}, "Reserved public IP address to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. Usage: --public-ips IP1,IP2")
	update.AddIntFlag(config.ArgLanId, "", 0, "The unique LAN Id of existing LANs to be attached to worker Nodes")
	update.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgK8sNodePoolId, config.ArgIdShort, "", config.RequiredFlagK8sNodePoolId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for NodePool/Request [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "nodepool",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Kubernetes NodePool",
		LongDesc: `This command deletes a Kubernetes Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`,
		Example:    deleteK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgK8sNodePoolId, config.ArgIdShort, "", config.RequiredFlagK8sNodePoolId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgK8sClusterId, config.ArgK8sNodePoolId)
}

func PreRunK8sClusterDcIdsNodePoolName(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgK8sClusterId, config.ArgDataCenterId, config.ArgName)
}

func RunK8sNodePoolList(c *core.CommandConfig) error {
	k8ss, _, err := c.K8s().ListNodePools(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePools(k8ss)))
}

func RunK8sNodePoolGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, GetStateK8sNodePool, viper.GetString(core.GetFlagName(c.NS, config.ArgK8sNodePoolId))); err != nil {
		return err
	}
	u, _, err := c.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolCreate(c *core.CommandConfig) error {
	newNodePool, err := getNewK8sNodePool(c)
	if err != nil {
		return err
	}
	u, _, err := c.K8s().CreateNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)), *newNodePool)
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, GetStateK8sNodePool, *id); err != nil {
				return err
			}
			if u, _, err = c.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)), *id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Node Pool id")
		}
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolUpdate(c *core.CommandConfig) error {
	oldNodePool, _, err := c.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	newNodePool := getNewK8sNodePoolUpdated(oldNodePool, c)
	_, _, err = c.K8s().UpdateNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgK8sNodePoolId)), newNodePool)
	if err != nil {
		return err
	}
	if err = utils.WaitForState(c, GetStateK8sNodePool, viper.GetString(core.GetFlagName(c.NS, config.ArgK8sNodePoolId))); err != nil {
		return err
	}
	//return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(newNodePoolUpdated))) (TO BE UPDATED)
	return c.Printer.Print(getK8sNodePoolPrint(c, nil))
}

func RunK8sNodePoolDelete(c *core.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s node pool")
	if err != nil {
		return err
	}
	_, err = c.K8s().DeleteNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	return c.Printer.Print("Status: Command node pool delete has been successfully executed")
}

// Wait for State

func GetStateK8sNodePool(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)), objId)
	if err != nil {
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}

func getNewK8sNodePool(c *core.CommandConfig) (*resources.K8sNodePool, error) {
	var (
		k8sversion string
		err        error
	)
	n := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, config.ArgK8sVersion))
	} else {
		if k8sversion, err = getK8sVersion(c); err != nil {
			return nil, err
		}
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	nodeCount := viper.GetInt32(core.GetFlagName(c.NS, config.ArgK8sNodeCount))
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, config.ArgCpuFamily))
	coresCount := viper.GetInt32(core.GetFlagName(c.NS, config.ArgCores))
	ramSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, config.ArgRam)), utils.MegaBytes)
	if err != nil {
		return nil, err
	}
	ram := int32(ramSize)
	nodeZone := viper.GetString(core.GetFlagName(c.NS, config.ArgAvailabilityZone))
	storageSize := viper.GetInt32(core.GetFlagName(c.NS, config.ArgStorageSize))
	storageType := viper.GetString(core.GetFlagName(c.NS, config.ArgStorageType))
	return &resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:             &n,
				K8sVersion:       &k8sversion,
				DatacenterId:     &dcId,
				NodeCount:        &nodeCount,
				CpuFamily:        &cpuFamily,
				CoresCount:       &coresCount,
				RamSize:          &ram,
				AvailabilityZone: &nodeZone,
				StorageType:      &storageType,
				StorageSize:      &storageSize,
			},
		},
	}, nil
}

func getNewK8sNodePoolUpdated(oldUser *resources.K8sNodePool, c *core.CommandConfig) resources.K8sNodePool {
	propertiesUpdated := resources.K8sNodePoolProperties{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sVersion)) {
			propertiesUpdated.SetK8sVersion(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sVersion)))
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sNodeCount)) {
			propertiesUpdated.SetNodeCount(viper.GetInt32(core.GetFlagName(c.NS, config.ArgK8sNodeCount)))
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMinNodeCount)) ||
			viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaxNodeCount)) {
			var minCount, maxCount int32
			autoScaling := properties.GetAutoScaling()
			if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMinNodeCount)) {
				minCount = viper.GetInt32(core.GetFlagName(c.NS, config.ArgK8sMinNodeCount))
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}
			if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaxNodeCount)) {
				maxCount = viper.GetInt32(core.GetFlagName(c.NS, config.ArgK8sMaxNodeCount))
			} else {
				if m, ok := autoScaling.GetMaxNodeCountOk(); ok && m != nil {
					maxCount = *m
				}
			}
			propertiesUpdated.SetAutoScaling(ionoscloud.KubernetesAutoScaling{
				MinNodeCount: &minCount,
				MaxNodeCount: &maxCount,
			})
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := getMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sAnnotationKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sAnnotationValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, config.ArgK8sAnnotationKey))
			value := viper.GetString(core.GetFlagName(c.NS, config.ArgK8sAnnotationValue))
			propertiesUpdated.SetAnnotations(map[string]string{
				key: value,
			})
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgLabelKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, config.ArgLabelValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
			value := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
			propertiesUpdated.SetLabels(map[string]string{
				key: value,
			})
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgLanId)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			lanId := viper.GetInt32(core.GetFlagName(c.NS, config.ArgLanId))
			newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{Id: &lanId})
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}
			propertiesUpdated.SetLans(newLans)
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgPublicIps)) {
			propertiesUpdated.SetPublicIps(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgPublicIps)))
		}
	}
	return resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &propertiesUpdated.KubernetesNodePoolProperties,
		},
	}
}

// Output Printing

var defaultK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "NodeCount", "DatacenterId", "State"}

var allK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "DatacenterId", "NodeCount", "CpuFamily", "StorageType", "State",
	"CoresCount", "RamSize", "AvailabilityZone", "StorageSize", "MaintenanceWindow", "AutoScaling", "PublicIps", "PublicIps", "AvailableUpgradeVersions"}

type K8sNodePoolPrint struct {
	NodePoolId               string   `json:"NodePoolId,omitempty"`
	Name                     string   `json:"Name,omitempty"`
	K8sVersion               string   `json:"K8sVersion,omitempty"`
	DatacenterId             string   `json:"DatacenterId,omitempty"`
	NodeCount                int32    `json:"NodeCount,omitempty"`
	CpuFamily                string   `json:"CpuFamily,omitempty"`
	StorageType              string   `json:"StorageType,omitempty"`
	State                    string   `json:"State,omitempty"`
	CoresCount               int32    `json:"CoresCount,omitempty"`
	RamSize                  int32    `json:"RamSize,omitempty"`
	AvailabilityZone         string   `json:"AvailabilityZone,omitempty"`
	StorageSize              int32    `json:"StorageSize,omitempty"`
	MaintenanceWindow        string   `json:"MaintenanceWindow,omitempty"`
	AutoScaling              string   `json:"AutoScaling,omitempty"`
	PublicIps                []string `json:"PublicIps,omitempty"`
	AvailableUpgradeVersions []string `json:"AvailableUpgradeVersions,omitempty"`
}

func getK8sNodePoolPrint(c *core.CommandConfig, k8ss []resources.K8sNodePool) printer.Result {
	r := printer.Result{}
	if c != nil {
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sNodePoolsKVMaps(k8ss)
			r.Columns = getK8sNodePoolCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getK8sNodePoolCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var k8sCols []string
		columnsMap := map[string]string{
			"NodePoolId":               "NodePoolId",
			"Name":                     "Name",
			"K8sVersion":               "K8sVersion",
			"DatacenterId":             "DatacenterId",
			"NodeCount":                "NodeCount",
			"CpuFamily":                "CpuFamily",
			"StorageType":              "StorageType",
			"State":                    "State",
			"CoresCount":               "CoresCount",
			"RamSize":                  "RamSize",
			"AvailabilityZone":         "AvailabilityZone",
			"StorageSize":              "StorageSize",
			"MaintenanceWindow":        "MaintenanceWindow",
			"AutoScaling":              "AutoScaling",
			"PublicIps":                "PublicIps",
			"AvailableUpgradeVersions": "AvailableUpgradeVersions",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				k8sCols = append(k8sCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return k8sCols
	} else {
		return defaultK8sNodePoolCols
	}
}

func getK8sNodePools(k8ss resources.K8sNodePools) []resources.K8sNodePool {
	u := make([]resources.K8sNodePool, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.K8sNodePool{KubernetesNodePool: item})
		}
	}
	return u
}

func getK8sNodePool(u *resources.K8sNodePool) []resources.K8sNodePool {
	k8ss := make([]resources.K8sNodePool, 0)
	if u != nil {
		k8ss = append(k8ss, resources.K8sNodePool{KubernetesNodePool: u.KubernetesNodePool})
	}
	return k8ss
}

func getK8sNodePoolsKVMaps(us []resources.K8sNodePool) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint K8sNodePoolPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.NodePoolId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				uPrint.Name = *name
			}
			if v, ok := properties.GetK8sVersionOk(); ok && v != nil {
				uPrint.K8sVersion = *v
			}
			if v, ok := properties.GetDatacenterIdOk(); ok && v != nil {
				uPrint.DatacenterId = *v
			}
			if v, ok := properties.GetNodeCountOk(); ok && v != nil {
				uPrint.NodeCount = *v
			}
			if v, ok := properties.GetCpuFamilyOk(); ok && v != nil {
				uPrint.CpuFamily = *v
			}
			if v, ok := properties.GetRamSizeOk(); ok && v != nil {
				uPrint.RamSize = *v
			}
			if v, ok := properties.GetStorageTypeOk(); ok && v != nil {
				uPrint.StorageType = *v
			}
			if v, ok := properties.GetStorageSizeOk(); ok && v != nil {
				uPrint.StorageSize = *v
			}
			if v, ok := properties.GetCoresCountOk(); ok && v != nil {
				uPrint.CoresCount = *v
			}
			if v, ok := properties.GetPublicIpsOk(); ok && v != nil {
				uPrint.PublicIps = *v
			}
			if v, ok := properties.GetAvailableUpgradeVersionsOk(); ok && v != nil {
				uPrint.AvailableUpgradeVersions = *v
			}
			if v, ok := properties.GetAvailabilityZoneOk(); ok && v != nil {
				uPrint.AvailabilityZone = *v
			}
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				if day, ok := maintenance.GetDayOfTheWeekOk(); ok && day != nil {
					uPrint.MaintenanceWindow = *day
				}
				if time, ok := maintenance.GetTimeOk(); ok && time != nil {
					uPrint.MaintenanceWindow = uPrint.MaintenanceWindow + " " + *time
				}
			}
			if autoScaling, ok := properties.GetAutoScalingOk(); ok && autoScaling != nil {
				if min, ok := autoScaling.GetMinNodeCountOk(); ok && min != nil {
					uPrint.AutoScaling = fmt.Sprintf("Min: %v", *min)
				}
				if max, ok := autoScaling.GetMaxNodeCountOk(); ok && max != nil {
					uPrint.AutoScaling = fmt.Sprintf("%s Max: %v", uPrint.AutoScaling, *max)
				}
			}
		}
		if meta, ok := u.GetMetadataOk(); ok && meta != nil {
			if state, ok := meta.GetStateOk(); ok && state != nil {
				uPrint.State = *state
			}
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}

func getK8sNodePoolsIds(outErr io.Writer, clusterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(clientSvc.Get(), context.TODO())
	k8ss, _, err := k8sSvc.ListNodePools(clusterId)
	clierror.CheckError(err, outErr)
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesNodePools.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				k8ssIds = append(k8ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return k8ssIds
}
