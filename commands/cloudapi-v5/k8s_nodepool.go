package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sNodePoolCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultK8sNodePoolCols, printer.ColsMessage(allK8sNodePoolCols))
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
		LongDesc:   "Use this command to get a list of all contained NodePools in a selected Kubernetes Cluster.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.K8sNodePoolsFiltersUsage() + "\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    listK8sNodePoolsExample,
		PreCmdRun:  PreRunK8sNodePoolsList,
		CmdRun:     RunK8sNodePoolList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsFilters(), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified NodePool to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "nodepool",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Kubernetes NodePool",
		LongDesc: `Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center. Regarding the name for the Kubernetes NodePool, the limit is 63 characters following the rule to begin and end with an alphanumeric character with dashes, underscores, dots, and alphanumerics between. Regarding the Kubernetes Version for the NodePool, if not set via flag, it will be used the default one: ` + "`" + `ionosctl k8s version get` + "`" + `.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run a command:

* K8s Cluster Id
* Datacenter Id`,
		Example:    createK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterDcIds,
		CmdRun:     RunK8sNodePoolCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "UnnamedNodePool", "The name for the K8s NodePool")
	create.AddStringFlag(cloudapiv5.ArgK8sVersion, "", "", "The K8s version for the NodePool. If not set, the default one will be used")
	create.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(cloudapiv5.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability")
	create.AddIntFlag(cloudapiv5.ArgCores, "", 2, "The total number of cores for the Node")
	create.AddStringFlag(cloudapiv5.ArgRam, "", strconv.Itoa(2048), "RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgCpuFamily, "", cloudapiv5.DefaultServerCPUFamily, "CPU Type")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgAvailabilityZone, cloudapiv5.ArgAvailabilityZoneShort, "AUTO", "The compute Availability Zone in which the Node should exist")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgStorageType, "", "HDD", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgStorageSize, "", strconv.Itoa(cloudapiv5.DefaultVolumeSize), "The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state[seconds]")

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
	update.AddStringFlag(cloudapiv5.ArgK8sVersion, "", "", "The K8s version for the NodePool. K8s version downgrade is not supported")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sVersion,
		func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			clusterId := viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgK8sClusterId))
			nodepoolId := viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgK8sNodePoolId))
			return completer.K8sNodePoolUpgradeVersions(os.Stderr, clusterId, nodepoolId), cobra.ShellCompDirectiveNoFileComp
		})
	update.AddIntFlag(cloudapiv5.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the NodePool should contain")
	update.AddIntFlag(cloudapiv5.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count")
	update.AddIntFlag(cloudapiv5.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count")
	update.AddStringFlag(cloudapiv5.ArgLabelKey, "", "", "Label key. Must be set together with --label-value")
	update.AddStringFlag(cloudapiv5.ArgLabelValue, "", "", "Label value. Must be set together with --label-key")
	update.AddStringFlag(cloudapiv5.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value")
	update.AddStringFlag(cloudapiv5.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key")
	update.AddStringFlag(cloudapiv5.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringSliceFlag(cloudapiv5.ArgPublicIps, "", []string{""}, "Reserved public IP address to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. Usage: --public-ips IP1,IP2")
	update.AddIntSliceFlag(cloudapiv5.ArgLanIds, "", []int{}, "The unique LAN Ids of existing LANs to be attached to worker Nodes. It will be attached to the existing ones")
	update.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

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
		PreCmdRun:  PreRunK8sClusterNodePoolDelete,
		CmdRun:     RunK8sNodePoolDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all the Kubernetes Node Pools within an existing Kubernetes Cluster.")

	return k8sCmd
}

func PreRunK8sNodePoolsList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgK8sClusterId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.K8sNodePoolsFilters(), completer.K8sNodePoolsFiltersUsage())
	}
	return nil
}

func PreRunK8sClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgK8sNodePoolId)
}

func PreRunK8sClusterNodePoolDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgK8sNodePoolId},
		[]string{cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgAll},
	)
}

func PreRunK8sClusterDcIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgDataCenterId)
}

func RunK8sNodePoolList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting K8s NodePools from K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters
			if val, ok := filters["ramSize"]; ok {
				convertedSize, err := utils.ConvertSize(val, utils.MegaBytes)
				if err != nil {
					return err
				}
				filters["ramSize"] = strconv.Itoa(convertedSize)
				listQueryParams.Filters = &filters
			}
		}
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	k8ss, resp, err := c.CloudApiV5Services.K8s().ListNodePools(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePools(k8ss)))
}

func RunK8sNodePoolGet(c *core.CommandConfig) error {
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	if err := utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, k8sNodePoolId); err != nil {
		return err
	}
	c.Printer.Verbose("K8s node pool with id: %v from K8s Cluster with id: %v is getting...", k8sNodePoolId, k8sClusterId)
	u, resp, err := c.CloudApiV5Services.K8s().GetNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	c.Printer.Verbose("Creating K8s NodePool in K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	u, resp, err := c.CloudApiV5Services.K8s().CreateNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)), *newNodePool)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, *id); err != nil {
				return err
			}
			if u, _, err = c.CloudApiV5Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)), *id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Node Pool id")
		}
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolUpdate(c *core.CommandConfig) error {
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	oldNodePool, _, err := c.CloudApiV5Services.K8s().GetNodePool(k8sClusterId, k8sNodePoolId)
	if err != nil {
		return err
	}
	newNodePool := getNewK8sNodePoolUpdated(oldNodePool, c)
	c.Printer.Verbose("Updating K8s node pool with id: %v from K8s Cluster with id: %v...", k8sNodePoolId, k8sClusterId)
	newNodePoolUpdated, resp, err := c.CloudApiV5Services.K8s().UpdateNodePool(k8sClusterId, k8sNodePoolId, newNodePool)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))); err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(newNodePoolUpdated)))
}

func RunK8sNodePoolDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	clusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	nodepollId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = DeleteAllK8sNodepools(c)
		if err != nil {
			return err
		}
	} else {
		err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s node pool")
		if err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting K8s Nodepool with id: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId)))
		resp, err = c.CloudApiV5Services.K8s().DeleteNodePool(clusterId, nodepollId)
		if resp != nil {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}
	return c.Printer.Print("Status: Command node pool delete has been successfully executed")
}

func getNewK8sNodePool(c *core.CommandConfig) (*resources.K8sNodePoolForPost, error) {
	var (
		k8sversion string
		err        error
	)
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion))
	} else {
		if k8sversion, err = getK8sVersion(c); err != nil {
			return nil, err
		}
	}
	ramSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgRam)), utils.MegaBytes)
	if err != nil {
		return nil, err
	}
	storageSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgStorageSize)), utils.GigaBytes)
	if err != nil {
		return nil, err
	}
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
	nodeCount := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeCount))
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgCpuFamily))
	cores := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgCores))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgAvailabilityZone))
	storageType := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgStorageType))
	// Set Properties
	nodePoolProperties := ionoscloud.KubernetesNodePoolPropertiesForPost{}
	nodePoolProperties.SetName(name)
	nodePoolProperties.SetK8sVersion(k8sversion)
	nodePoolProperties.SetNodeCount(nodeCount)
	nodePoolProperties.SetDatacenterId(dcId)
	nodePoolProperties.SetCpuFamily(cpuFamily)
	nodePoolProperties.SetCoresCount(cores)
	nodePoolProperties.SetRamSize(int32(ramSize))
	nodePoolProperties.SetAvailabilityZone(availabilityZone)
	nodePoolProperties.SetStorageSize(int32(storageSize))
	nodePoolProperties.SetStorageType(storageType)

	c.Printer.Verbose("Properties set for creating the node pool: Name: %v, K8sVersion: %v, NodeCount: %v, DatacenterId: %v, CpuFamily: %v, CoresCount: %v, RamSize: %vMB, AvailabilityZone: %v, StorageSize: %v, StorageType: %v",
		name, k8sversion, nodeCount, dcId, cpuFamily, cores, int32(ramSize), availabilityZone, int32(storageSize), storageType)

	return &resources.K8sNodePoolForPost{
		KubernetesNodePoolForPost: ionoscloud.KubernetesNodePoolForPost{
			Properties: &nodePoolProperties,
		},
	}, nil
}

func getNewK8sNodePoolUpdated(oldUser *resources.K8sNodePool, c *core.CommandConfig) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion)) {
			version := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion))
			propertiesUpdated.SetK8sVersion(version)
			c.Printer.Verbose("Property K8sVersion set: %v", version)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeCount)) {
			nodeCount := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeCount))
			propertiesUpdated.SetNodeCount(nodeCount)
			c.Printer.Verbose("Property NodeCount set: %v", nodeCount)
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMinNodeCount)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaxNodeCount)) {
			var minCount, maxCount int32
			autoScaling := properties.GetAutoScaling()
			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMinNodeCount)) {
				minCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMinNodeCount))
				c.Printer.Verbose("Property MinNodeCount set: %v", minCount)
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}
			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaxNodeCount)) {
				maxCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaxNodeCount))
				c.Printer.Verbose("Property MaxNodeCount set: %v", maxCount)
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
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := getMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
				c.Printer.Verbose("Property MaintenanceWindow set: %v", newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sAnnotationKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sAnnotationValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sAnnotationKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sAnnotationValue))
			propertiesUpdated.SetAnnotations(map[string]string{
				key: value,
			})
			c.Printer.Verbose("Property Annotations set: key: %v, value: %v", key, value)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
			propertiesUpdated.SetLabels(map[string]string{
				key: value,
			})
			c.Printer.Verbose("Property Labels set: key: %v, value: %v", key, value)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgLanIds)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}
			// Add new LANs
			lanIds := viper.GetIntSlice(core.GetFlagName(c.NS, cloudapiv5.ArgLanIds))
			for _, lanId := range lanIds {
				id := int32(lanId)
				newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
					Id: &id,
				})
				c.Printer.Verbose("Property Lans set: %v", id)
			}
			propertiesUpdated.SetLans(newLans)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgPublicIps)) {
			publicIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv5.ArgPublicIps))
			propertiesUpdated.SetPublicIps(publicIps)
			c.Printer.Verbose("Property PublicIps set: %v", publicIps)
		}
	}
	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

func DeleteAllK8sNodepools(c *core.CommandConfig) (*resources.Response, error) {
	clusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	_ = c.Printer.Print("K8sNodePools to be deleted:")
	k8sNodePools, resp, err := c.CloudApiV5Services.K8s().ListNodePools(clusterId, resources.ListQueryParams{})
	if err != nil {
		return nil, err
	}
	if k8sNodePoolsItems, ok := k8sNodePools.GetItemsOk(); ok && k8sNodePoolsItems != nil {
		for _, dc := range *k8sNodePoolsItems {
			if id, ok := dc.GetIdOk(); ok && id != nil {
				_ = c.Printer.Print("K8sNodePool Id: " + *id)
			}
			if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print("K8sNodePool Name: " + *name)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the K8sNodePools"); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the K8sNodePools")

		for _, dc := range *k8sNodePoolsItems {
			if id, ok := dc.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Starting deleting K8sNodePool with id: %v...", *id)
				resp, err = c.CloudApiV5Services.K8s().DeleteNodePool(clusterId, *id)
				if resp != nil {
					c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
				}
				if err != nil {
					return nil, err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return nil, err
				}
			}
		}
	}
	return resp, nil
}

// Output Printing

var defaultK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "NodeCount", "DatacenterId", "State"}

var allK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "DatacenterId", "NodeCount", "CpuFamily", "StorageType", "State", "LanIds",
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
	LanIds                   []int32  `json:"LanIds,omitempty"`
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
			"LanIds":                   "LanIds",
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
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				lanIds := make([]int32, 0)
				for _, lanItem := range *lans {
					if lanId, ok := lanItem.GetIdOk(); ok && lanId != nil {
						lanIds = append(lanIds, *lanId)
					}
				}
				uPrint.LanIds = lanIds
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
