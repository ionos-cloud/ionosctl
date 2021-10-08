package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
		LongDesc:   "Use this command to get a list of all contained NodePools in a selected Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    listK8sNodePoolsExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sNodePoolList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgK8sNodePoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified NodePool to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

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

Note: If you want to attach multiple LANs to Node Pool, use ` + "`" + `--lan-ids=LAN_ID1,LAN_ID2` + "`" + ` flag. If you want to also set a Route Network, Route GatewayIp for LAN use ` + "`" + `ionosctl k8s nodepool lan add` + "`" + ` command for each LAN you want to add.

Required values to run a command:

* K8s Cluster Id
* Datacenter Id`,
		Example:    createK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterDcIds,
		CmdRun:     RunK8sNodePoolCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "UnnamedNodePool", "The name for the K8s NodePool")
	create.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "The K8s version for the NodePool. If not set, the default one will be used")
	create.AddStringFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "Collection of LAN Ids of existing LANs to be attached to worker Nodes")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanIds, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	create.AddIntFlag(cloudapiv6.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability")
	create.AddIntFlag(cloudapiv6.ArgCores, "", 2, "The total number of cores for the Node")
	create.AddStringFlag(cloudapiv6.ArgRam, "", strconv.Itoa(2048), "RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgCpuFamily, "", cloudapiv6.DefaultServerCPUFamily, "CPU Type")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgAvailabilityZone, cloudapiv6.ArgAvailabilityZoneShort, "AUTO", "The compute Availability Zone in which the Node should exist")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgStorageType, "", "HDD", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgStorageSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state[seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "nodepool",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Kubernetes NodePool",
		LongDesc: `Use this command to update the number of worker Nodes, the minimum and maximum number of worker Nodes, the add labels, annotations, to update the maintenance day and time, to attach private LANs to a Node Pool within an existing Kubernetes Cluster. You can also add reserved public IP addresses to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. The array must contain one extra IP than maximum number of Nodes could be. The extra provided IP Will be used during rebuilding of Nodes.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Note: If you want to attach multiple LANs to Node Pool, use ` + "`" + `--lan-ids=LAN_ID1,LAN_ID2` + "`" + ` flag. If you want to also set a Route Network, Route GatewayIp for LAN use ` + "`" + `ionosctl k8s nodepool lan add` + "`" + ` command for each LAN you want to add.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`,
		Example:    updateK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "The K8s version for the NodePool. K8s version downgrade is not supported")
	update.AddIntFlag(cloudapiv6.ArgK8sNodeCount, "", 1, "The number of worker Nodes that the NodePool should contain")
	update.AddIntFlag(cloudapiv6.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count")
	update.AddIntFlag(cloudapiv6.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count")
	update.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", "Label key. Must be set together with --label-value")
	update.AddStringFlag(cloudapiv6.ArgLabelValue, "", "", "Label value. Must be set together with --label-key")
	update.AddStringFlag(cloudapiv6.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value")
	update.AddStringFlag(cloudapiv6.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key")
	update.AddStringFlag(cloudapiv6.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringSliceFlag(cloudapiv6.ArgPublicIps, "", []string{""}, "Reserved public IP address to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. Usage: --public-ips IP1,IP2")
	update.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "Collection of LAN Ids of existing LANs to be attached to worker Nodes. It will be added to the existing LANs attached")
	update.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	update.AddStringFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgK8sNodePoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

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
	deleteCmd.AddStringFlag(cloudapiv6.ArgK8sClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgK8sNodePoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	k8sCmd.AddCommand(K8sNodePoolLanCmd())

	return k8sCmd
}

func PreRunK8sClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgK8sClusterId, cloudapiv6.ArgK8sNodePoolId)
}

func PreRunK8sClusterDcIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgK8sClusterId, cloudapiv6.ArgDataCenterId)
}

func RunK8sNodePoolList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting K8s NodePools from K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)))
	k8ss, resp, err := c.CloudApiV6Services.K8s().ListNodePools(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePools(k8ss)))
}

func RunK8sNodePoolGet(c *core.CommandConfig) error {
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId))
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId))
	if err := utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, k8sNodePoolId); err != nil {
		return err
	}
	c.Printer.Verbose("K8s node pool with id: %v from K8s Cluster with id: %v is getting...", k8sNodePoolId, k8sClusterId)
	u, resp, err := c.CloudApiV6Services.K8s().GetNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
	c.Printer.Verbose("Creating K8s NodePool in K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)))
	u, resp, err := c.CloudApiV6Services.K8s().CreateNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)), *newNodePool)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, *id); err != nil {
				return err
			}
			if u, _, err = c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)), *id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Node Pool id")
		}
	}
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(u)))
}

func RunK8sNodePoolUpdate(c *core.CommandConfig) error {
	oldNodePool, _, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId)))
	if err != nil {
		return err
	}
	newNodePool := getNewK8sNodePoolUpdated(oldNodePool, c)
	_, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId)), newNodePool)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId))); err != nil {
		return err
	}
	newNodePoolUpdated, _, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId)))
	return c.Printer.Print(getK8sNodePoolPrint(c, getK8sNodePool(newNodePoolUpdated)))
}

func RunK8sNodePoolDelete(c *core.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s node pool")
	if err != nil {
		return err
	}
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodePoolId))
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sClusterId))
	c.Printer.Verbose("Deleting K8s node pool with id: %v from K8s Cluster with id: %v...", k8sNodePoolId, k8sClusterId)
	resp, err := c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print("Status: Command node pool delete has been successfully executed")
}

func getNewK8sNodePool(c *core.CommandConfig) (*resources.K8sNodePool, error) {
	var (
		k8sversion string
		err        error
	)
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion))
	} else {
		if k8sversion, err = getK8sVersion(c); err != nil {
			return nil, err
		}
	}
	ramSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRam)), utils.MegaBytes)
	if err != nil {
		return nil, err
	}
	storageSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgStorageSize)), utils.GigaBytes)
	if err != nil {
		return nil, err
	}
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	nodeCount := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeCount))
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCpuFamily))
	cores := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCores))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAvailabilityZone))
	storageType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgStorageType))
	// Set Properties
	nodePoolProperties := ionoscloud.KubernetesNodePoolProperties{}
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
	// Add LANs
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds)) {
		newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
		lanIds := viper.GetIntSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds))
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
		for _, lanId := range lanIds {
			id := int32(lanId)
			newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
				Id:   &id,
				Dhcp: &dhcp,
			})
		}
		nodePoolProperties.SetLans(newLans)
	}
	return &resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &nodePoolProperties,
		},
	}, nil
}

func getNewK8sNodePoolUpdated(oldUser *resources.K8sNodePool, c *core.CommandConfig) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion)) {
			vers := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion))
			propertiesUpdated.SetK8sVersion(vers)
			c.Printer.Verbose("Property K8sVersion set: %v", vers)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeCount)) {
			nodeCount := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeCount))
			propertiesUpdated.SetNodeCount(nodeCount)
			c.Printer.Verbose("Property NodeCount set: %v", nodeCount)
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMinNodeCount)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount)) {
			var minCount, maxCount int32
			autoScaling := properties.GetAutoScaling()
			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMinNodeCount)) {
				minCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMinNodeCount))
				c.Printer.Verbose("Property MinNodeCount set: %v", minCount)
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}
			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount)) {
				maxCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount))
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
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := getMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationValue))
			propertiesUpdated.SetAnnotations(map[string]string{
				key: value,
			})
			c.Printer.Verbose("Property Annotations set: key: %v, value: %v", key, value)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
			propertiesUpdated.SetLabels(map[string]string{
				key: value,
			})
			c.Printer.Verbose("Property Labels set: key: %v, value: %v", key, value)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			// Append existing LANs
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}
			// Add new LANs
			lanIds := viper.GetIntSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds))
			dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
			for _, lanId := range lanIds {
				id := int32(lanId)
				newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
					Id:   &id,
					Dhcp: &dhcp,
				})
				c.Printer.Verbose("Property Lans set: %v", id)
			}
			propertiesUpdated.SetLans(newLans)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublicIps)) {
			publicIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPublicIps))
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
