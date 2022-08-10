package dataplatform

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/commands/dataplatform/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-autoscaling"
)

func NodePoolCmd() *core.Command {
	ctx := context.TODO()
	nodePoolCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nodepool",
			Aliases:          []string{"np"},
			Short:            "Data Platform NodePool Operations",
			Long:             "The sub-commands of `ionosctl dataplatform nodepool` allow you to list, get, create, update, delete Data Platform NodePools.",
			TraverseChildren: true,
		},
	}
	globalFlags := nodePoolCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultNodePoolCols, printer.ColsMessage(allNodePoolCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(nodePoolCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = nodePoolCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNodePoolCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, nodePoolCmd, core.CommandBuilder{
		Namespace:  "dataplatform",
		Resource:   "nodepool",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Data Platform NodePools",
		LongDesc:   "Use this command to get a list of all contained NodePools in a selected Data Platform Cluster.\n\nRequired values to run command:\n\n*  Cluster Id",
		Example:    listNodePoolsExample,
		PreCmdRun:  PreRunNodePoolsList,
		CmdRun:     RunNodePoolList,
		InitClient: true,
	})
	list.AddStringFlag(dp.ArgClusterId, "", "", dp.ClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, nodePoolCmd, core.CommandBuilder{
		Namespace:  "dataplatform",
		Resource:   "nodepool",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Data Platform NodePool",
		LongDesc:   "Use this command to retrieve details about a specific NodePool from an existing Data Platform Cluster. You can wait for the Node Pool to be in \"ACTIVE\" state using `--wait-for-state` flag together with `--timeout` option.\n\nRequired values to run command:\n\n*  Cluster Id\n*  NodePool Id",
		Example:    getNodePoolExample,
		PreCmdRun:  PreRunClusterNodePoolIds,
		CmdRun:     RunNodePoolGet,
		InitClient: true,
	})
	get.AddStringFlag(dp.ArgClusterId, "", "", dp.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(dp.ArgNodePoolId, dp.ArgIdShort, "", dp.NodePoolId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(dp.ArgNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, dp.ArgClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified NodePool to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.TimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, nodePoolCmd, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Data Platform NodePool",
		LongDesc: `Use this command to create a Node Pool into an existing Data Platform Cluster. The Data Platform Cluster must be in state "ACTIVE" before creating a Node Pool. 

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run a command:

* Cluster Id
* Name
* Node Count
`,
		Example:    createNodePoolExample,
		PreCmdRun:  PreRunClusterNodePoolCreate,
		CmdRun:     RunNodePoolCreate,
		InitClient: true,
	})
	create.AddStringFlag(dp.ArgName, dp.ArgNameShort, "", "The name for the  NodePool", core.RequiredFlagOption())
	create.AddIntFlag(dp.ArgNodeCount, "", 0, "The number of nodes that make up the node pool.", core.RequiredFlagOption())
	create.AddStringFlag(dp.ArgClusterId, "", "", dp.ClusterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgCpuFamily, "", dp.DefaultServerCPUFamily, "A valid CPU family name or `AUTO` if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource.")
	create.AddIntFlag(dp.ArgCores, "", 0, "The number of CPU cores per node.")
	create.AddStringFlag(dp.ArgRam, "", "", "The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgAvailabilityZone, dp.ArgAvailabilityZoneShort, "AUTO", "The availability zone of the virtual datacenter region where the node pool resources should be provisioned.")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgStorageType, "", "", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgStorageSize, "", "", "The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgMaintenanceTime, dp.ArgMaintenanceTimeShort, "", "Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	create.AddStringFlag(dp.ArgMaintenanceDay, dp.ArgMaintenanceDayShort, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringToStringFlag(dp.ArgLabels, dp.ArgLabelsShort, map[string]string{}, "Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/). Use the following format: --labels KEY=VALUE,KEY=VALUE")
	create.AddStringToStringFlag(dp.ArgAnnotations, dp.ArgAnnotationsShort, map[string]string{}, "Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/. Use the following format: --annotations KEY=VALUE,KEY=VALUE")
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.TimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state[seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, nodePoolCmd, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Data Platform NodePool",
		LongDesc: `Use this command to update the number of worker Nodes, to add labels, annotations, to update the maintenance day and time. 

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.
Required values to run command:

*  Cluster Id
*  NodePool Id`,
		Example:    updateNodePoolExample,
		PreCmdRun:  PreRunClusterNodePoolIds,
		CmdRun:     RunNodePoolUpdate,
		InitClient: true,
	})
	update.AddStringFlag(dp.ArgClusterId, "", "", dp.ClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dp.ArgNodePoolId, dp.ArgIdShort, "", dp.NodePoolId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(dp.ArgNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, dp.ArgClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(dp.ArgNodeCount, "", 0, "The number of nodes that make up the node pool.")
	update.AddStringFlag(dp.ArgMaintenanceTime, dp.ArgMaintenanceTimeShort, "", "Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	update.AddStringFlag(dp.ArgMaintenanceDay, dp.ArgMaintenanceDayShort, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	_ = update.Command.RegisterFlagCompletionFunc(dp.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringToStringFlag(dp.ArgLabels, dp.ArgLabelsShort, map[string]string{}, "Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/). Use the following format: --labels KEY=VALUE,KEY=VALUE")
	update.AddStringToStringFlag(dp.ArgAnnotations, dp.ArgAnnotationsShort, map[string]string{}, "Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/. Use the following format: --annotations KEY=VALUE,KEY=VALUE")
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.TimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state[seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, nodePoolCmd, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Data Platform NodePool",
		LongDesc: `This command deletes a Data Platform Node Pool within an existing Data Platform Cluster.

Required values to run command:

*  Cluster Id
*  NodePool Id`,
		Example:    deleteNodePoolExample,
		PreCmdRun:  PreRunClusterNodePoolDelete,
		CmdRun:     RunNodePoolDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(dp.ArgClusterId, "", "", dp.ClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(dp.ArgNodePoolId, dp.ArgIdShort, "", dp.NodePoolId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(dp.ArgNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, dp.ArgClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgAll, config.ArgAllShort, false, "Delete all the Data Platform Node Pools within an existing Data Platform Cluster.")

	return nodePoolCmd
}

func PreRunNodePoolsList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, dp.ArgClusterId); err != nil {
		return err
	}
	return nil
}

func PreRunClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, dp.ArgClusterId, dp.ArgNodePoolId)
}

func PreRunClusterNodePoolCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, dp.ArgClusterId, dp.ArgName, dp.ArgNodeCount)
	if err != nil {
		return err
	} // Validate Flags
	if len(viper.GetString(core.GetFlagName(c.NS, dp.ArgName))) > 63 {
		return errors.New("name string has to have a maximum of 63 characters")
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgCores)) {
		if viper.GetInt(core.GetFlagName(c.NS, dp.ArgCores)) < 1 {
			return errors.New("cores count should have a value of at least 1")
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgRam)) {
		if viper.GetInt(core.GetFlagName(c.NS, dp.ArgRam)) < 2048 {
			return errors.New("ram size should have a value of at least 2048")
		}
		if viper.GetInt(core.GetFlagName(c.NS, dp.ArgRam))%1024 != 0 {
			return errors.New("ram size should be a multiple of 1024")
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgStorageSize)) {
		if viper.GetInt(core.GetFlagName(c.NS, dp.ArgStorageSize)) < 10 {
			return errors.New("storage size should have a value of at least 10")
		}
	}
	return nil

}

func PreRunClusterNodePoolDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{dp.ArgClusterId, dp.ArgNodePoolId},
		[]string{dp.ArgClusterId, config.ArgAll},
	)
}

func RunNodePoolList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting  NodePools from  Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))
	nodepools, resp, err := c.DataPlatformServices.NodePools().List(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))

	if err != nil {
		return err
	}
	return c.Printer.Print(getNodePoolPrint(resp, c, getNodePools(nodepools)))
}

func RunNodePoolGet(c *core.CommandConfig) error {
	NodePoolId := viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId))
	ClusterId := viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))
	if err := utils.WaitForState(c, waiter.NodePoolStateInterrogator, NodePoolId); err != nil {
		return err
	}
	c.Printer.Verbose(" node pool with id: %v from  Cluster with id: %v is getting...", NodePoolId, ClusterId)
	nodePool, _, err := c.DataPlatformServices.NodePools().Get(ClusterId, NodePoolId)

	if err != nil {
		return err
	}
	return c.Printer.Print(getNodePoolPrint(nil, c, getNodePool(&nodePool)))
}

func RunNodePoolCreate(c *core.CommandConfig) error {
	newNodePool, err := getNewNodePool(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Creating  NodePool in  Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))
	nodePool, resp, err := c.DataPlatformServices.NodePools().Create(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)), *newNodePool)
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := nodePool.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.NodePoolStateInterrogator, *id); err != nil {
				return err
			}
			if nodePool, _, err = c.DataPlatformServices.NodePools().Get(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)), *id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new  Node Pool id")
		}
	}
	return c.Printer.Print(getNodePoolPrint(resp, c, []resources.NodePoolResponseData{nodePool}))
}

func RunNodePoolUpdate(c *core.CommandConfig) error {
	_ = c.Printer.Print("WARNING: The following flags are deprecated:" + c.Command.GetAnnotationsByKey(core.DeprecatedFlagsAnnotation) + ". Use --labels, --annotations options instead!")
	oldNodePool, _, err := c.DataPlatformServices.NodePools().Get(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)),
		viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId)))
	if err != nil {
		return err
	}
	newNodePool := getNewNodePoolUpdated(&oldNodePool, c)
	_, resp, err := c.DataPlatformServices.NodePools().Update(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)),
		viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId)), newNodePool)
	if err != nil {
		return err
	}
	if err = utils.WaitForState(c, waiter.NodePoolStateInterrogator, viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId))); err != nil {
		return err
	}
	newNodePoolUpdated, _, err := c.DataPlatformServices.NodePools().Get(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)),
		viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId)))
	return c.Printer.Print(getNodePoolPrint(resp, c, getNodePool(&newNodePoolUpdated)))
}

func RunNodePoolDelete(c *core.CommandConfig) error {
	ClusterId := viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))
	NodePoolId := viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId))
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgAll)) {
		if err := DeleteAllNodepools(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		err := utils.AskForConfirm(c.Stdin, c.Printer, "delete  node pool")
		if err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting  node pool with id: %v from  Cluster with id: %v...", NodePoolId, ClusterId)
		_, _, err = c.DataPlatformServices.NodePools().Delete(ClusterId, NodePoolId)
		if err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	}
}

func getNewNodePool(c *core.CommandConfig) (*resources.CreateNodePoolRequest, error) {
	name := viper.GetString(core.GetFlagName(c.NS, dp.ArgName))
	nodeCount := viper.GetInt32(core.GetFlagName(c.NS, dp.ArgNodeCount))

	// Set Properties
	nodePoolProperties := ionoscloud.CreateNodePoolProperties{}
	nodePoolProperties.SetName(name)
	c.Printer.Verbose("Property Name set: %v", name)
	nodePoolProperties.SetNodeCount(nodeCount)
	c.Printer.Verbose("Property NodeCount set: %v", nodeCount)
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgCpuFamily)) {
		cpuFamily := viper.GetString(core.GetFlagName(c.NS, dp.ArgCpuFamily))
		nodePoolProperties.SetCpuFamily(cpuFamily)
		c.Printer.Verbose("Property CPU Family set: %v", cpuFamily)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, dp.ArgCores))
		nodePoolProperties.SetCoresCount(cores)
		c.Printer.Verbose("Property CoresCount set: %v", cores)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgRam)) {
		ramSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, dp.ArgRam)), utils.MegaBytes)
		if err != nil {
			return nil, err
		}
		nodePoolProperties.SetRamSize(int32(ramSize))
		c.Printer.Verbose("Property RAM Size set: %vMB", int32(ramSize))
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgAvailabilityZone)) {
		availabilityZone := viper.GetString(core.GetFlagName(c.NS, dp.ArgAvailabilityZone))
		nodePoolProperties.SetAvailabilityZone(ionoscloud.AvailabilityZone(availabilityZone))
		c.Printer.Verbose("Property Availability Zone set: %v", availabilityZone)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgAvailabilityZone)) {
		storageType := viper.GetString(core.GetFlagName(c.NS, dp.ArgStorageType))
		nodePoolProperties.SetStorageType(ionoscloud.StorageType(storageType))
		c.Printer.Verbose("Property Storage Type set: %v", storageType)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgAvailabilityZone)) {
		storageSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, dp.ArgStorageSize)), utils.GigaBytes)
		if err != nil {
			return nil, err
		}
		nodePoolProperties.SetStorageSize(int32(storageSize))
		c.Printer.Verbose("Property Storage Size set: %vGB", int32(storageSize))

	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgLabels)) {
		keyValueMapLabels := viper.GetStringMap(core.GetFlagName(c.NS, dp.ArgLabels))
		nodePoolProperties.SetLabels(keyValueMapLabels)
		c.Printer.Verbose("Property Labels set: %v", keyValueMapLabels)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgAnnotations)) {
		keyValueMapAnnotations := viper.GetStringMap(core.GetFlagName(c.NS, dp.ArgAnnotations))
		nodePoolProperties.SetAnnotations(keyValueMapAnnotations)
		c.Printer.Verbose("Property Annotations set: %v", keyValueMapAnnotations)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
		maintenanceWindow := ionoscloud.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceTime))
			c.Printer.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceDay))
			c.Printer.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(maintenanceDay)
		}
		nodePoolProperties.SetMaintenanceWindow(maintenanceWindow)
	}
	return &resources.CreateNodePoolRequest{
		CreateNodePoolRequest: ionoscloud.CreateNodePoolRequest{
			Properties: &nodePoolProperties,
		},
	}, nil
}

func getNewNodePoolUpdated(oldNodePool *resources.NodePoolResponseData, c *core.CommandConfig) resources.PatchNodePoolRequest {
	propertiesUpdated := resources.PatchNodePoolProperties{}
	if properties, ok := oldNodePool.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgNodeCount)) {
			nodeCount := viper.GetInt32(core.GetFlagName(c.NS, dp.ArgNodeCount))
			propertiesUpdated.SetNodeCount(nodeCount)
			c.Printer.Verbose("Property NodeCount set: %v", nodeCount)
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) {
			newMaintenanceWindow := ionoscloud.MaintenanceWindow{}
			if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) {
				maintenanceTime := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceTime))
				c.Printer.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
				newMaintenanceWindow.SetTime(maintenanceTime)
			}
			if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
				maintenanceDay := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceDay))
				c.Printer.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
				newMaintenanceWindow.SetDayOfTheWeek(maintenanceDay)
			}
			propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgLabels)) {
			keyValueMapLabels := viper.GetStringMap(core.GetFlagName(c.NS, dp.ArgLabels))
			propertiesUpdated.SetLabels(keyValueMapLabels)
			c.Printer.Verbose("Property Labels set: %v", keyValueMapLabels)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgAnnotations)) {
			keyValueMapAnnotations := viper.GetStringMap(core.GetFlagName(c.NS, dp.ArgAnnotations))
			propertiesUpdated.SetAnnotations(keyValueMapAnnotations)
			c.Printer.Verbose("Property Annotations set: %v", keyValueMapAnnotations)
		}
	}
	return resources.PatchNodePoolRequest{
		PatchNodePoolRequest: ionoscloud.PatchNodePoolRequest{
			Properties: &propertiesUpdated.PatchNodePoolProperties,
		},
	}
}

func DeleteAllNodepools(c *core.CommandConfig) error {
	ClusterId := viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))
	c.Printer.Verbose("Cluster ID: %v", ClusterId)
	c.Printer.Verbose("Getting NodePools...")
	NodePools, _, err := c.DataPlatformServices.NodePools().List(ClusterId)
	if err != nil {
		return err
	}
	if NodePoolsItems, ok := NodePools.GetItemsOk(); ok && NodePoolsItems != nil {
		if len(*NodePoolsItems) > 0 {
			_ = c.Printer.Print("NodePools to be deleted:")
			for _, dc := range *NodePoolsItems {
				toPrint := ""
				if id, ok := dc.GetIdOk(); ok && id != nil {
					toPrint += "NodePool Id: " + *id
				}
				if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += " NodePool Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the NodePools"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the NodePools")
			var multiErr error
			for _, dc := range *NodePoolsItems {
				if id, ok := dc.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting NodePool with id: %v...", *id)
					_, _, err = c.DataPlatformServices.NodePools().Delete(ClusterId, *id)
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForDelete(c, waiter.NodePoolDeleteInterrogator, *id); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.WaitDeleteAllAppendErr, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no NodePools found")
		}
	} else {
		return errors.New("could not get items of NodePools")
	}
}

// Output Printing

var defaultNodePoolCols = []string{"NodePoolId", "Name", "Version", "NodeCount", "DatacenterId", "State"}

var allNodePoolCols = []string{"NodePoolId", "Name", "Version", "DatacenterId", "NodeCount", "CpuFamily", "StorageType", "State", "LanIds",
	"CoresCount", "RamSize", "AvailabilityZone", "StorageSize", "MaintenanceWindow", "AutoScaling", "PublicIps", "AvailableUpgradeVersions",
	"Annotations", "Labels"}

type NodePoolPrint struct {
	NodePoolId        string                 `json:"NodePoolId,omitempty"`
	Name              string                 `json:"Name,omitempty"`
	NodeCount         int32                  `json:"NodeCount,omitempty"`
	CpuFamily         string                 `json:"CpuFamily,omitempty"`
	CoresCount        int32                  `json:"CoresCount,omitempty"`
	RamSize           int32                  `json:"RamSize,omitempty"`
	AvailabilityZone  string                 `json:"AvailabilityZone,omitempty"`
	StorageType       string                 `json:"StorageType,omitempty"`
	StorageSize       int32                  `json:"StorageSize,omitempty"`
	State             string                 `json:"State,omitempty"`
	MaintenanceWindow string                 `json:"MaintenanceWindow,omitempty"`
	Annotations       map[string]interface{} `json:"Annotations,omitempty"`
	Labels            map[string]interface{} `json:"Labels,omitempty"`
	Version           string                 `json:"Version,omitempty"`
	DatacenterId      string                 `json:"DatacenterId,omitempty"`
}

func getNodePoolPrint(resp *resources.Response, c *core.CommandConfig, nrd []resources.NodePoolResponseData) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if nrd != nil {
			r.OutputJSON = nrd
			r.KeyValue = getNodePoolsKVMaps(nrd)
			r.Columns = getNodePoolCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getNodePoolCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var Cols []string
		columnsMap := map[string]string{
			"NodePoolId":        "NodePoolId",
			"Name":              "Name",
			"NodeCount":         "NodeCount",
			"CpuFamily":         "CpuFamily",
			"CoresCount":        "CoresCount",
			"RamSize":           "RamSize",
			"AvailabilityZone":  "AvailabilityZone",
			"StorageType":       "StorageType",
			"StorageSize":       "StorageSize",
			"State":             "State",
			"MaintenanceWindow": "MaintenanceWindow",
			"Annotations":       "Annotations",
			"Labels":            "Labels",
			"Version":           "Version",
			"DatacenterId":      "DatacenterId",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				Cols = append(Cols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return Cols
	} else {
		return defaultNodePoolCols
	}
}

func getNodePools(s resources.NodePoolListResponseData) []resources.NodePoolResponseData {
	u := make([]resources.NodePoolResponseData, 0)
	if items, ok := s.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.NodePoolResponseData{NodePoolResponseData: item})
		}
	}
	return u
}

func getNodePool(u *resources.NodePoolResponseData) []resources.NodePoolResponseData {
	s := make([]resources.NodePoolResponseData, 0)
	if u != nil {
		s = append(s, resources.NodePoolResponseData{NodePoolResponseData: u.NodePoolResponseData})
	}
	return s
}

func getNodePoolsKVMaps(us []resources.NodePoolResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint NodePoolPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.NodePoolId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if nameOk, ok := properties.GetNameOk(); ok && nameOk != nil {
				uPrint.Name = *nameOk
			}
			if nodeCountOk, ok := properties.GetNodeCountOk(); ok && nodeCountOk != nil {
				uPrint.NodeCount = *nodeCountOk
			}
			if cpuFamilyOk, ok := properties.GetCpuFamilyOk(); ok && cpuFamilyOk != nil {
				uPrint.CpuFamily = *cpuFamilyOk
			}
			if coresCountOk, ok := properties.GetCoresCountOk(); ok && coresCountOk != nil {
				uPrint.CoresCount = *coresCountOk
			}
			if ramSizeOk, ok := properties.GetRamSizeOk(); ok && ramSizeOk != nil {
				uPrint.RamSize = *ramSizeOk
			}
			if availabilityZoneOk, ok := properties.GetAvailabilityZoneOk(); ok && availabilityZoneOk != nil {
				uPrint.AvailabilityZone = string(*availabilityZoneOk)
			}
			if storageTypeOk, ok := properties.GetStorageTypeOk(); ok && storageTypeOk != nil {
				uPrint.StorageType = string(*storageTypeOk)
			}
			if storageSizeOk, ok := properties.GetStorageSizeOk(); ok && storageSizeOk != nil {
				uPrint.StorageSize = *storageSizeOk
			}
			if maintenanceWindowOk, ok := properties.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				if dayOfTheWeekOk, ok := maintenanceWindowOk.GetDayOfTheWeekOk(); ok && dayOfTheWeekOk != nil {
					uPrint.MaintenanceWindow = *dayOfTheWeekOk
				}
				if timeOk, ok := maintenanceWindowOk.GetTimeOk(); ok && timeOk != nil {
					uPrint.MaintenanceWindow = uPrint.MaintenanceWindow + " " + *timeOk
				}
			}
			if annotationsOk, ok := properties.GetAnnotationsOk(); ok && annotationsOk != nil {
				uPrint.Annotations = *annotationsOk
			}
			if labelsOk, ok := properties.GetLabelsOk(); ok && labelsOk != nil {
				uPrint.Labels = *labelsOk
			}

			if versionOk, ok := properties.GetDataPlatformVersionOk(); ok && versionOk != nil {
				uPrint.Version = *versionOk
			}
			if datacenterIdOk, ok := properties.GetDatacenterIdOk(); ok && datacenterIdOk != nil {
				uPrint.DatacenterId = *datacenterIdOk
			}
		}
		if metadataOk, ok := u.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				uPrint.State = *stateOk
			}
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}
