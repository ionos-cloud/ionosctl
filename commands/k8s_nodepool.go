package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8sNodePool() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "nodepool",
			Short:            "Kubernetes NodePool Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s nodepool` + "`" + ` allow you to list, get, create, update, delete Kubernetes NodePools.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultK8sNodePoolCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(k8sCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sNodePoolCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterId, RunK8sNodePoolList, "list", "List Kubernetes NodePools",
		"Use this command to get a list of all contained NodePools in a selected Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id", listK8sNodePoolsExample, true)
	list.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIds, RunK8sNodePoolGet, "get", "Get a Kubernetes NodePool",
		"Use this command to retrieve details about a specific NodePool from an existing Kubernetes Cluster. You can wait for the Node Pool to be in \"ACTIVE\" state using `+\"`\"+`--wait-for-state`+\"`\"+` flag together with `+\"`\"+`--timeout`+\"`\"+` option.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id",
		getK8sNodePoolExample, true)
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), get.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, "", config.DefaultWait, "Wait for specified NodePool to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, "", config.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterDcIdsNodePoolName, RunK8sNodePoolCreate, "create", "Create a Kubernetes NodePool",
		`Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center. Regarding the name for the Kubernetes NodePool, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.

You can wait for the Node Pool to be in "ACTIVE" state using `+"`"+`--wait-for-state`+"`"+` flag together with `+"`"+`--timeout`+"`"+` option.

Required values to run a command:

* K8s Cluster Id
* Datacenter Id
* K8s NodePool Name`, createK8sNodePoolExample, true)
	create.AddStringFlag(config.ArgK8sNodePoolName, "", "", "The name for the K8s NodePool "+config.RequiredFlag)
	create.AddStringFlag(config.ArgK8sNodePoolVersion, "", "", "The K8s version for the NodePool")
	create.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability")
	create.AddIntFlag(config.ArgCoresCount, "", 2, "The total number of cores for the Node")
	create.AddIntFlag(config.ArgRamSize, "", 2048, "The amount of memory for the node in MB, e.g. 2048. Size must be specified in multiples of 1024 MB (1 GB) with a minimum of 2048 MB")
	create.AddStringFlag(config.ArgCpuFamily, "", config.DefaultServerCPUFamily, "CPU Type")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgK8sNodeZone, "", "AUTO", "The compute Availability Zone in which the Node should exist")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgStorageType, "", "HDD", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgStorageSize, "", 10, "The total allocated storage capacity of a Node")
	create.AddBoolFlag(config.ArgWaitForState, "", config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, "", config.K8sTimeoutSeconds, "Timeout option for waiting for NodePool/Request [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIds, RunK8sNodePoolUpdate, "update", "Update a Kubernetes NodePool",
		`Use this command to update the number of worker Nodes, the minimum and maximum number of worker Nodes, the add labels, annotations, to update the maintenance day and time, to attach private LANs to a Node Pool within an existing Kubernetes Cluster.

You can wait for the Node Pool to be in "ACTIVE" state using `+"`"+`--wait-for-state`+"`"+` flag together with `+"`"+`--timeout`+"`"+` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`, updateK8sNodePoolExample, true)
	update.AddStringFlag(config.ArgK8sNodePoolVersion, "", "", "The K8s version for the NodePool. K8s version downgrade is not supported")
	update.AddIntFlag(config.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the NodePool should contain")
	update.AddIntFlag(config.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count")
	update.AddIntFlag(config.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count")
	update.AddStringFlag(config.ArgLabelKey, "", "", "Label key. Must be set together with --label-value")
	update.AddStringFlag(config.ArgLabelValue, "", "", "Label value. Must be set together with --label-key")
	update.AddStringFlag(config.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value")
	update.AddStringFlag(config.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key")
	update.AddStringFlag(config.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	update.AddStringFlag(config.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddIntFlag(config.ArgLanId, "", 0, "The unique LAN Id of existing LANs to be attached to worker Nodes")
	update.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), update.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, "", config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	update.AddIntFlag(config.ArgTimeout, "", config.K8sTimeoutSeconds, "Timeout option for waiting for NodePool/Request [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIds, RunK8sNodePoolDelete, "delete", "Delete a Kubernetes NodePool",
		`This command deletes a Kubernetes Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`, deleteK8sNodePoolExample, true)
	deleteCmd.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), deleteCmd.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sClusterNodePoolIds(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId, config.ArgK8sNodePoolId)
}

func PreRunK8sClusterDcIdsNodePoolName(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId, config.ArgDataCenterId, config.ArgK8sNodePoolName)
}

func RunK8sNodePoolList(c *builder.CommandConfig) error {
	k8ss, _, err := c.K8s().ListNodePools(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePools(k8ss)))
}

func RunK8sNodePoolGet(c *builder.CommandConfig) error {
	if viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForState)) {
		if err := utils.WaitForState(c, GetStateK8sNodePool, viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId))); err != nil {
			return err
		}
	}
	u, _, err := c.K8s().GetNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolCreate(c *builder.CommandConfig) error {
	newNodePool := getNewK8sNodePool(c)
	u, _, err := c.K8s().CreateNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)), newNodePool)
	if err != nil {
		return err
	}
	if viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, GetStateK8sNodePool, *id); err != nil {
				return err
			}
			if u, _, err = c.K8s().GetNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)), *id); err != nil {
				return err
			}
		}
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolUpdate(c *builder.CommandConfig) error {
	oldNodePool, _, err := c.K8s().GetNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	newNodePool := getNewK8sNodePoolUpdated(oldNodePool, c)
	newNodePoolUpdated, _, err := c.K8s().UpdateNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)), newNodePool)
	if err != nil {
		return err
	}
	if viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForState)) {
		if err = utils.WaitForState(c, GetStateK8sNodePool, viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getK8sNodePoolUpdatedPrint(c, newNodePoolUpdated))
}

func RunK8sNodePoolDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s node pool")
	if err != nil {
		return err
	}
	_, err = c.K8s().DeleteNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	return c.Printer.Print("Status: Command node pool delete has been successfully executed")
}

// Wait for State

func GetStateK8sNodePool(c *builder.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.K8s().GetNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)), objId)
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

func getNewK8sNodePool(c *builder.CommandConfig) resources.K8sNodePool {
	n := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolName))
	v := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolVersion))
	dcId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId))
	nodeCount := viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeCount))
	cpuFamily := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgCpuFamily))
	coresCount := viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgCoresCount))
	ramSize := viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgRamSize))
	nodeZone := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeZone))
	storageSize := viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgStorageSize))
	storageType := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgStorageType))
	return resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:             &n,
				K8sVersion:       &v,
				DatacenterId:     &dcId,
				NodeCount:        &nodeCount,
				CpuFamily:        &cpuFamily,
				CoresCount:       &coresCount,
				RamSize:          &ramSize,
				AvailabilityZone: &nodeZone,
				StorageType:      &storageType,
				StorageSize:      &storageSize,
			},
		},
	}
}

func getNewK8sNodePoolUpdated(oldUser *resources.K8sNodePool, c *builder.CommandConfig) resources.K8sNodePool {
	propertiesUpdated := resources.K8sNodePoolProperties{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolVersion)) {
			propertiesUpdated.SetK8sVersion(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolVersion)))
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeCount)) {
			propertiesUpdated.SetNodeCount(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeCount)))
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMinNodeCount)) ||
			viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaxNodeCount)) {
			var minCount, maxCount int32
			autoScaling := properties.GetAutoScaling()
			if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMinNodeCount)) {
				minCount = viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMinNodeCount))
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}
			if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaxNodeCount)) {
				maxCount = viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaxNodeCount))
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
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceDay)) ||
			viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := getMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sAnnotationKey)) &&
			viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sAnnotationValue)) {
			key := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sAnnotationKey))
			value := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sAnnotationValue))
			propertiesUpdated.SetAnnotations(ionoscloud.KubernetesNodePoolAnnotation{
				Key:   &key,
				Value: &value,
			})
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)) &&
			viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue)) {
			key := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey))
			value := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue))
			propertiesUpdated.SetLabels(ionoscloud.KubernetesNodePoolLabel{
				Key:   &key,
				Value: &value,
			})
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			lanId := viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId))
			newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{Id: &lanId})
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}
			propertiesUpdated.SetLans(newLans)
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

func getK8sNodePoolPrint(c *builder.CommandConfig, k8ss []resources.K8sNodePool) printer.Result {
	r := printer.Result{}
	if c != nil {
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sNodePoolsKVMaps(k8ss)
			r.Columns = getK8sNodePoolCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getK8sNodePoolUpdatedPrint(c *builder.CommandConfig, k8ss *resources.K8sNodePoolUpdated) printer.Result {
	r := printer.Result{}
	if c != nil {
		r.Resource = c.ParentName
		r.Verb = c.Name
		if k8ss != nil {
			r.OutputJSON = k8ss
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
