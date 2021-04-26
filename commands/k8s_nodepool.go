package commands

import (
	"context"
	"errors"
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
			Use:              "k8s-nodepool",
			Short:            "K8s NodePool Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s-nodepool` + "`" + ` allow you to list, get, create, update, delete K8s NodePools.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultK8sNodePoolCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(k8sCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	list := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterIdValidate, RunK8sNodePoolList, "list", "List K8s NodePools",
		"Use this command to get a list of existing K8s NodePools.\n\nRequired values to run command:\n\n* K8s Cluster Id", listK8sNodePoolsExample, true)
	list.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIdsValidate, RunK8sNodePoolGet, "get", "Get a K8s NodePool",
		"Use this command to retrieve details about a specific K8s NodePool.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id",
		getK8sNodePoolExample, true)
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Command.Name(), get.Command.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	create := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterIdNodePoolNameValidate, RunK8sNodePoolCreate, "create", "Create a K8s NodePool",
		`Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center.

Required values to run a command:

* K8s NodePool Name`, createK8sNodePoolExample, true)
	create.AddStringFlag(config.ArgK8sNodePoolName, "", "", "The name for the K8s NodePool "+config.RequiredFlag)
	create.AddStringFlag(config.ArgK8sNodePoolVersion, "", "1.19.8", "The K8s version for the NodePool")
	create.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgK8sNodeCount, "", 2, "The number of worker Nodes that the Node Pool should contain. Min 2, Max: Determined by the resource availability")
	create.AddIntFlag(config.ArgCoresCount, "", 2, "The total number of cores for the Node")
	create.AddIntFlag(config.ArgRamSize, "", 2048, "The amount of memory for the node in MB, e.g. 2048. Size must be specified in multiples of 1024 MB (1 GB) with a minimum of 2048 MB")
	create.AddStringFlag(config.ArgCpuFamily, "", config.DefaultServerCPUFamily, "Cpu Type")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgK8sNodeZone, "", "AUTO", "The compute availability zone in which the node should exist")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgStorageType, "", "HDD", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgStorageSize, "", 10, "The total allocated storage capacity of a Node")

	update := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIdsValidate, RunK8sNodePoolUpdate, "update", "Update a K8s NodePool",
		`Use this command to update the number of worker Nodes or other properties for a Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`, updateK8sNodePoolExample, true)
	update.AddStringFlag(config.ArgK8sNodePoolVersion, "", "", "The K8s version for the NodePool. K8s version downgrade is not supported")
	update.AddStringFlag(config.ArgK8sNodeCount, "", "", "The number of worker Nodes that the NodePool should contain")
	update.AddIntFlag(config.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count")
	update.AddIntFlag(config.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count")
	update.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Command.Name(), update.Command.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	deleteCmd := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIdsValidate, RunK8sNodePoolDelete, "delete", "Delete a K8s NodePool",
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
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Command.Name(), deleteCmd.Command.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sClusterNodePoolIdsValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId, config.ArgK8sNodePoolId)
}

func PreRunK8sClusterIdNodePoolNameValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId, config.ArgK8sNodePoolName)
}

func RunK8sNodePoolList(c *builder.CommandConfig) error {
	k8ss, _, err := c.K8s().ListNodePools(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePools(k8ss)))
}

func RunK8sNodePoolGet(c *builder.CommandConfig) error {
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
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolUpdate(c *builder.CommandConfig) error {
	oldNodePool, _, err := c.K8s().GetNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	_, _, err = c.K8s().UpdateNodePool(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)), getNewK8sNodePoolUpdated(oldNodePool, c))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, nil))
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
	return c.Printer.Print(getK8sNodePoolPrint(c, nil))
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
	var minCount, maxCount int32
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
	return resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &propertiesUpdated.KubernetesNodePoolProperties,
		},
	}
}

// Output Printing

var defaultK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "DatacenterId", "NodeCount", "CpuFamily", "StorageType", "State"}

type K8sNodePoolPrint struct {
	NodePoolId   string `json:"NodePoolId,omitempty"`
	Name         string `json:"Name,omitempty"`
	K8sVersion   string `json:"K8sVersion,omitempty"`
	DatacenterId string `json:"DatacenterId,omitempty"`
	NodeCount    int32  `json:"NodeCount,omitempty"`
	CpuFamily    string `json:"CpuFamily,omitempty"`
	StorageType  string `json:"StorageType,omitempty"`
	State        string `json:"State,omitempty"`
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

func getK8sNodePoolCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var k8sCols []string
		columnsMap := map[string]string{
			"NodePoolId":   "NodePoolId",
			"Name":         "Name",
			"K8sVersion":   "K8sVersion",
			"DatacenterId": "DatacenterId",
			"NodeCount":    "NodeCount",
			"CpuFamily":    "CpuFamily",
			"StorageType":  "StorageType",
			"State":        "State",
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
