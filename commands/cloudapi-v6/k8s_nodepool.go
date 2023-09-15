package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

var (
	allK8sNodepoolJSONPaths = map[string]string{
		"NodePoolId":               "id",
		"Name":                     "properties.name",
		"K8sVersion":               "properties.k8sVersion",
		"DatacenterId":             "properties.datacenterId",
		"NodeCount":                "properties.nodeCount",
		"CpuFamily":                "properties.cpuArchitecture",
		"StorageType":              "properties.storageType",
		"State":                    "metadata.state",
		"LanIds":                   "properties.lan.*.id",
		"CoresCount":               "properties.coresCount",
		"RamSize":                  "properties.ramSize",
		"AvailabilityZone":         "properties.availabilityZone",
		"StorageSize":              "properties.storageSize",
		"MaintenanceWindow":        "properties.maintenanceWindow",
		"AutoScaling":              "properties.autoScaling",
		"PublicIps":                "properties.publicIps",
		"AvailableUpgradeVersions": "properties.availableUpgradeVersions",
		"Annotations":              "properties.annotations",
		"Labels":                   "properties.labels",
	}

	defaultK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "NodeCount", "DatacenterId", "State"}
	allK8sNodePoolCols     = []string{"NodePoolId", "Name", "K8sVersion", "DatacenterId", "NodeCount", "CpuFamily", "StorageType", "State", "LanIds",
		"CoresCount", "RamSize", "AvailabilityZone", "StorageSize", "MaintenanceWindow", "AutoScaling", "PublicIps", "AvailableUpgradeVersions",
		"Annotations", "Labels", "ClusterId"}
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultK8sNodePoolCols, printer.ColsMessage(allK8sNodePoolCols))
	_ = viper.BindPFlag(core.GetFlagName(k8sCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

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
	get.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(constants.FlagNodepoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for specified NodePool to be in ACTIVE state")
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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

Required values to run a command (for Public Kubernetes Cluster):

* K8s Cluster Id
* Datacenter Id

Required values to run a command (for Private Kubernetes Cluster):

* K8s Cluster Id
* Datacenter Id`,
		Example:    createK8sNodePoolExample,
		PreCmdRun:  PreRunK8sClusterDcIds,
		CmdRun:     RunK8sNodePoolCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "UnnamedNodePool", "The name for the K8s NodePool")
	create.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "The K8s version for the NodePool. If not set, the default one will be used")
	create.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "Collection of LAN Ids of existing LANs to be attached to worker Nodes")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanIds, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	create.AddIntFlag(constants.FlagNodeCount, "", 1, "The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability")
	create.AddIntFlag(constants.FlagCores, "", 2, "The total number of cores for the Node")
	create.AddStringFlag(constants.FlagRam, "", strconv.Itoa(2048), "RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCpuFamily, "", cloudapiv6.DefaultServerCPUFamily, "CPU Type")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(create.Command.Context(), os.Stderr, datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "The compute Availability Zone in which the Node should exist")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageType, "", "HDD", "Storage Type")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	create.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")
	create.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state[seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sVersion,
		func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			clusterId := viper.GetString(core.GetFlagName(update.NS, constants.FlagClusterId))
			nodepoolId := viper.GetString(core.GetFlagName(update.NS, constants.FlagNodepoolId))
			return completer.K8sNodePoolUpgradeVersions(os.Stderr, clusterId, nodepoolId), cobra.ShellCompDirectiveNoFileComp
		})
	update.AddIntFlag(constants.FlagNodeCount, "", 1, "The number of worker Nodes that the NodePool should contain")
	update.AddIntFlag(cloudapiv6.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count")
	update.AddIntFlag(cloudapiv6.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count")
	update.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	update.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")
	update.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", "Label key. Must be set together with --label-value", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	update.AddStringFlag(cloudapiv6.ArgLabelValue, "", "", "Label value. Must be set together with --label-key", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	update.AddStringFlag(cloudapiv6.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	update.AddStringFlag(cloudapiv6.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	update.AddStringFlag(cloudapiv6.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringSliceFlag(cloudapiv6.ArgPublicIps, "", []string{}, "Reserved public IP address to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. Usage: --public-ips IP1,IP2")
	update.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "Collection of LAN Ids of existing LANs to be attached to worker Nodes. It will be added to the existing LANs attached")
	update.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	update.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(constants.FlagNodepoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
	deleteCmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(constants.FlagNodepoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Kubernetes Node Pools within an existing Kubernetes Nodepools.")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	k8sCmd.AddCommand(K8sNodePoolLanCmd())

	return k8sCmd
}

func PreRunK8sNodePoolsList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.K8sNodePoolsFilters(), completer.K8sNodePoolsFiltersUsage())
	}
	return nil
}

func PreRunK8sClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId)
}

func PreRunK8sClusterNodePoolDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId, constants.FlagNodepoolId},
		[]string{constants.FlagClusterId, cloudapiv6.ArgAll},
	)
}

func PreRunK8sClusterDcIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, constants.FlagClusterId},
		[]string{cloudapiv6.ArgDataCenterId, constants.FlagClusterId})
}

func RunK8sNodePoolListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	clusters, _, err := c.CloudApiV6Services.K8s().ListClusters(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	var allNodePools []ionoscloud.KubernetesNodePools
	var allNodePoolsConverted []map[string]interface{}
	totalTime := time.Duration(0)

	for _, cluster := range getK8sClusters(clusters) {
		nodePools, resp, err := c.CloudApiV6Services.K8s().ListNodePools(*cluster.GetId(), listQueryParams)
		if err != nil {
			return err
		}

		items, ok := nodePools.GetItemsOk()
		if !ok || items == nil {
			continue
		}

		for _, node := range *items {
			temp, err := json2table.ConvertJSONToTable("", allK8sNodepoolJSONPaths, node)
			if err != nil {
				return fmt.Errorf("failed to convert from JSON to Table format: %w", err)
			}

			clusterId, ok := cluster.GetIdOk()
			if !ok || clusterId == nil {
				continue
			}

			temp[0]["ClusterId"] = clusterId
			allNodePoolsConverted = append(allNodePoolsConverted, temp[0])
		}

		allNodePools = append(allNodePools, nodePools.KubernetesNodePools)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(allNodePools, allNodePoolsConverted,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunK8sNodePoolList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunK8sNodePoolListAll(c)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting K8s NodePools from K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters

			if val, ok := filters["ramSize"]; ok {
				convertedSize, err := utils.ConvertSize(val[0], utils.MegaBytes)
				if err != nil {
					return err
				}

				filters["ramSize"] = []string{strconv.Itoa(convertedSize)}
				listQueryParams.Filters = &filters
			}
		}
	}

	k8ss, resp, err := c.CloudApiV6Services.K8s().ListNodePools(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allK8sNodepoolJSONPaths, k8ss.KubernetesNodePools,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunK8sNodePoolGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	if err := utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, k8sNodePoolId); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"K8s node pool with id: %v from K8s Cluster with id: %v is getting...", k8sNodePoolId, k8sClusterId))

	u, resp, err := c.CloudApiV6Services.K8s().GetNodePool(k8sClusterId, k8sNodePoolId, queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allK8sNodepoolJSONPaths, u.KubernetesNodePool,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunK8sNodePoolCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	newNodePool, err := getNewK8sNodePool(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Creating K8s NodePool in K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	u, resp, err := c.CloudApiV6Services.K8s().CreateNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), *newNodePool, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, *id); err != nil {
				return err
			}
			if u, _, err = c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), *id, queryParams); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Node Pool id")
		}
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allK8sNodepoolJSONPaths, u.KubernetesNodePool,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunK8sNodePoolUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	oldNodePool, _, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)), queryParams)
	if err != nil {
		return err
	}

	newNodePool := getNewK8sNodePoolUpdated(oldNodePool, c)
	_, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)), newNodePool, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForState(c, waiter.K8sNodePoolStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))); err != nil {
		return err
	}

	newNodePoolUpdated, _, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)), queryParams)

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allK8sNodepoolJSONPaths, newNodePoolUpdated.KubernetesNodePool,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunK8sNodePoolDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err = DeleteAllK8sNodepools(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Kubernetes Nodepools successfully deleted"))
		return nil
	}

	if !confirm.Ask("delete k8s node pool") {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting K8s node pool with id: %v from K8s Cluster with id: %v...", k8sNodePoolId, k8sClusterId))

	resp, err := c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, k8sNodePoolId, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Kubernetes Nodepool successfully deleted"))
	return nil
}

func getNewK8sNodePool(c *core.CommandConfig) (*resources.K8sNodePoolForPost, error) {
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

	ramSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils.MegaBytes)
	if err != nil {
		return nil, err
	}

	storageSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils.GigaBytes)
	if err != nil {
		return nil, err
	}

	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	nodeCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagNodeCount))
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily))
	cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))
	storageType := viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageType))

	// Set Properties
	nodePoolProperties := ionoscloud.KubernetesNodePoolPropertiesForPost{}
	nodePoolProperties.SetName(name)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))

	nodePoolProperties.SetK8sVersion(k8sversion)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property K8sVersion set: %v", k8sversion))

	nodePoolProperties.SetNodeCount(nodeCount)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NodeCount set: %v", nodeCount))

	nodePoolProperties.SetDatacenterId(dcId)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DatacenterId set: %v", dcId))

	nodePoolProperties.SetCpuFamily(cpuFamily)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CPU Family set: %v", cpuFamily))

	nodePoolProperties.SetCoresCount(cores)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CoresCount set: %v", cores))

	nodePoolProperties.SetRamSize(int32(ramSize))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property RAM Size set: %vMB", int32(ramSize)))

	nodePoolProperties.SetAvailabilityZone(availabilityZone)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Availability Zone set: %v", availabilityZone))

	nodePoolProperties.SetStorageSize(int32(storageSize))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Storage Size set: %vGB", int32(storageSize)))

	nodePoolProperties.SetStorageType(storageType)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Storage Type set: %v", storageType))

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLabels)) {
		keyValueMapLabels := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagLabels))
		nodePoolProperties.SetLabels(keyValueMapLabels)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Labels set: %v", keyValueMapLabels))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagAnnotations)) {
		keyValueMapAnnotations := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagAnnotations))
		nodePoolProperties.SetAnnotations(keyValueMapAnnotations)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Annotations set: %v", keyValueMapAnnotations))
	}

	// Add LANs
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds)) {
		newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
		lanIds := viper.GetIntSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds))
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))

		for _, lanId := range lanIds {
			id := int32(lanId)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Lan ID set: %v", id))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Dhcp set: %v", dhcp))

			newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
				Id:   &id,
				Dhcp: &dhcp,
			})
		}

		nodePoolProperties.SetLans(newLans)
	}

	return &resources.K8sNodePoolForPost{
		KubernetesNodePoolForPost: ionoscloud.KubernetesNodePoolForPost{
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

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property K8sVersion set: %v", vers))
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagNodeCount)) {
			nodeCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagNodeCount))
			propertiesUpdated.SetNodeCount(nodeCount)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NodeCount set: %v", nodeCount))
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

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property MinNodeCount set: %v", minCount))
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}

			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount)) {
				maxCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount))

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property MaxNodeCount set: %v", maxCount))
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

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Annotations set: key: %v, value: %v", key, value))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
			propertiesUpdated.SetLabels(map[string]string{
				key: value,
			})

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Labels set: key: %v, value: %v", key, value))
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLabels)) {
			keyValueMapLabels := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagLabels))
			propertiesUpdated.SetLabels(keyValueMapLabels)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Labels set: %v", keyValueMapLabels))
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagAnnotations)) {
			keyValueMapAnnotations := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagAnnotations))
			propertiesUpdated.SetAnnotations(keyValueMapAnnotations)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Annotations set: %v", keyValueMapAnnotations))
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

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Lans set: %v", id))
			}

			propertiesUpdated.SetLans(newLans)
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublicIps)) {
			publicIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPublicIps))
			propertiesUpdated.SetPublicIps(publicIps)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property PublicIps set: %v", publicIps))
		}
	}

	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

func DeleteAllK8sNodepools(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("K8sCluster ID: %v", k8sClusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting K8sNodePools..."))

	k8sNodePools, resp, err := c.CloudApiV6Services.K8s().ListNodePools(k8sClusterId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	k8sNodePoolsItems, ok := k8sNodePools.GetItemsOk()
	if !ok || k8sNodePoolsItems == nil {
		return fmt.Errorf("could not get items of Kubernetes Nodepools")
	}

	if len(*k8sNodePoolsItems) <= 0 {
		return fmt.Errorf("no Kubernetes Nodepools found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("K8sNodePools to be deleted:"))

	for _, dc := range *k8sNodePoolsItems {
		delIdAndName := ""

		if id, ok := dc.GetIdOk(); ok && id != nil {
			delIdAndName += "K8sNodePool Id: " + *id
		}

		if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " K8sNodePool Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the K8sNodePools") {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the K8sNodePools"))

	var multiErr error
	for _, dc := range *k8sNodePoolsItems {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting K8sNodePool with id: %v...", *id))

		resp, err = c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, *id, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}

		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
			continue
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
