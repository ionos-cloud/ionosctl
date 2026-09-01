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
		ShortDesc: "Create a Kubernetes NodePool (worker Nodes)",
		LongDesc: `Create a node pool of worker Nodes inside an existing Managed Kubernetes cluster.
The cluster must be in state ACTIVE first. Every Node in the pool is provisioned
into the Data Center named by --datacenter-id and shares the same hardware shape
(--cores, --ram, --storage-size, --storage-type, --cpu-family) and Kubernetes
version.

Name: up to 63 characters, must begin and end with an alphanumeric character,
with dashes, underscores, dots and alphanumerics in between.

Kubernetes version: if --k8s-version is not set, the parent cluster's version is
used. The pool version must be less than or equal to the cluster version.

CPU family: if --cpu-family is omitted, the first CPU family available in the
Data Center's location is chosen automatically. --server-type (VCPU or
DedicatedCore) selects the compute engine server type for the Nodes.

Autoscaling and reserved public IPs are not configurable at create time via
dedicated flags; set them afterwards with ` + "`ionosctl compute k8s nodepool update`" + `,
or pass a full JSON body with --json-properties.

Networking: attach existing LANs with --lan-ids=LAN_ID1,LAN_ID2. To also set a
route network and gateway IP on a LAN, use ` + "`ionosctl compute k8s nodepool lan add`" + `
per LAN after creation.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the node pool reaches the AVAILABLE state.

Required values to run command:

* K8s Cluster Id
* Datacenter Id`,
		Example: `# Create a minimal node pool (1 Node, defaults for hardware and version) and wait for it
ionosctl compute k8s nodepool create --cluster-id CLUSTER_ID --datacenter-id DATACENTER_ID --name pool-a --wait

# Create a 3-Node pool with a specific shape, SSD storage, pinned Kubernetes version and two attached LANs
ionosctl compute k8s nodepool create --cluster-id CLUSTER_ID --datacenter-id DATACENTER_ID --name workers \
  --node-count 3 --cores 4 --ram 8GB --storage-type SSD --storage-size 100 --k8s-version 1.29.5 --lan-ids 1,2`,
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
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "UnnamedNodePool", "Name of the node pool. Up to 63 characters; must start and end with an alphanumeric character, dashes, underscores, dots and alphanumerics in between")
	cmd.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "Kubernetes version for the worker Nodes, e.g. 1.29.5. Must be <= the cluster's version. If not set, the cluster's version is used")
	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "IDs of existing LANs (in the same Data Center) to attach to the worker Nodes, e.g. --lan-ids 1,2. Use `nodepool lan add` to also set routes on a LAN")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanIds, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Whether Nodes obtain an IP on the attached LANs via DHCP. Applies to the LANs given in --lan-ids. e.g. --dhcp=true, --dhcp=false")
	cmd.AddIntFlag(constants.FlagNodeCount, "", 1, "Number of worker Nodes in the pool. Minimum 1; the maximum depends on your contract and resource availability")
	cmd.AddIntFlag(constants.FlagCores, "", 2, "Number of CPU cores per Node")
	cmd.AddStringFlag(constants.FlagRam, "", strconv.Itoa(2048), "RAM per Node. Minimum 2048MB and must be a multiple of 1024MB. Accepts a unit suffix, e.g. --ram 2048, --ram 2048MB or --ram 4GB")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "3GB", "4GB", "5GB", "10GB", "50GB", "100GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagCpuFamily, "", "",
		"CPU family for the Nodes (e.g. INTEL_SKYLAKE, INTEL_XEON), constrained by the Data Center's location. "+
			"If not set, the first CPU family available in that location (as returned by the API) is used")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(cmd.Command.Context(), datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "Compute availability zone for the Nodes. AUTO lets IONOS place them; ZONE_1 / ZONE_2 pin them to a specific zone")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagServerType, "", "", []string{"DedicatedCore", "VCPU"},
		"Compute-engine server type for the Nodes: 'DedicatedCore' (dedicated physical CPU cores) "+
			"or 'VCPU' (shared vCPU cores, typically cheaper)")
	cmd.AddStringFlag(constants.FlagStorageType, "", "HDD", "Type of the per-Node boot storage: HDD (default) or SSD")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "Per-Node boot storage size in GB. Accepts a unit suffix, e.g. --storage-size 10 or --storage-size 10GB. The maximum is bounded by your contract limit")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Kubernetes labels propagated onto the pool's Nodes (usable for scheduling/nodeSelectors). Overwrites any existing labels. Format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Kubernetes annotations propagated onto the pool's Nodes. Overwrites any existing annotations. Format: --annotations KEY=VALUE,KEY=VALUE")

	return cmd
}
