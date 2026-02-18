package nodepool

import (
	"context"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sNodePoolCreateCmd() *core.Command {
	jsonPropertiesExample := "{\n  \"metadata\": {},\n  \"properties\": {\n    \"name\": \"K8s-node-pool\",\n    \"datacenterId\": \"12345678-90ab-cdef-1234-567890abcdef\",\n    \"nodeCount\": 2,\n    \"cpuFamily\": \"INTEL_SKYLAKE\",\n    \"coresCount\": 4,\n    \"ramSize\": 2048,\n    \"availabilityZone\": \"AUTO\",\n    \"storageType\": \"HDD\",\n    \"storageSize\": 100,\n    \"k8sVersion\": \"1.27.6\",\n    \"maintenanceWindow\": {\n      \"dayOfTheWeek\": \"Monday\",\n      \"time\": \"13:00:00\"\n    },\n    \"autoScaling\": {\n      \"minNodeCount\": \"1\",\n      \"maxNodeCount\": \"2\"\n    },\n    \"lans\": [\n      {\n        \"id\": 1,\n        \"dhcp\": true,\n        \"routes\": [\n          {\n            \"network\": \"1.2.3.4/24\",\n            \"gatewayIp\": \"10.1.5.16\"\n          }\n        ]\n      }\n    ],\n    \"labels\": {\n      \"property1\": \"string\",\n      \"property2\": \"string\"\n    },\n    \"annotations\": {\n      \"property1\": \"string\",\n      \"property2\": \"string\"\n    },\n    \"publicIps\": [\n      \"203.0.113.1\",\n      \"203.0.113.2\",\n      \"203.0.113.3\"\n    ]\n  }\n}"
	nodepoolViaJsonPropertiesFlag := ionoscloud.KubernetesNodePoolForPost{}
	cmd := core.NewCommandWithJsonProperties(context.TODO(), nil, jsonPropertiesExample, &nodepoolViaJsonPropertiesFlag, core.CommandBuilder{
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
		Example: "ionosctl k8s nodepool create --cluster-id CLUSTER_ID --datacenter-id DATACENTER_ID --name NAME",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{cloudapiv6.ArgDataCenterId, constants.FlagClusterId},
				[]string{cloudapiv6.ArgDataCenterId, constants.FlagClusterId},
				[]string{constants.FlagJsonProperties, constants.FlagClusterId},
				[]string{constants.FlagJsonPropertiesExample},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if c.Command.Command.Flags().Changed(constants.FlagJsonProperties) {
				return RunK8sNodePoolCreateFromJSON(c, nodepoolViaJsonPropertiesFlag)
			}

			return RunK8sNodePoolCreate(c)
		},
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "UnnamedNodePool", "The name for the K8s NodePool")
	cmd.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "The K8s version for the NodePool. If not set, the version of the cluster will be used")
	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "Collection of LAN Ids of existing LANs to be attached to worker Nodes")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanIds, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	cmd.AddIntFlag(constants.FlagNodeCount, "", 1, "The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability")
	cmd.AddIntFlag(constants.FlagCores, "", 2, "The total number of cores for the Node")
	cmd.AddStringFlag(constants.FlagRam, "", strconv.Itoa(2048), "RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagCpuFamily, "", "",
		"CPU Type. If the flag is not set, the CPU Family will be chosen based on the location of the Datacenter. "+
			"It will always be the first CPU Family available, as returned by the API")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(cmd.Command.Context(), datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "The compute Availability Zone in which the Node should exist")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagServerType, "", "", []string{"DedicatedCore", "VCPU"},
		"The type of server for the Kubernetes node pool can be either"+
			"'DedicatedCore' (nodes with dedicated CPU cores) or 'VCPU' (nodes with shared CPU cores)."+
			"This selection corresponds to the server type for the compute engine.")
	cmd.AddStringFlag(constants.FlagStorageType, "", "HDD", "Storage Type")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")
	cmd.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state[seconds]")

	return cmd
}
