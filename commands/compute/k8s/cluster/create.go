package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func K8sClusterCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Kubernetes Cluster (control plane)",
		LongDesc: `Create a new Managed Kubernetes cluster - the control plane only. The cluster
starts with no worker capacity; add worker Nodes afterwards with
` + "`ionosctl compute k8s nodepool create`" + ` once the cluster reaches ACTIVE.

Name: up to 63 characters, must begin and end with an alphanumeric character,
with dashes, underscores, dots and alphanumerics in between.

Kubernetes version: if --k8s-version is not set, the current default is used
(see ` + "`ionosctl compute k8s version get`" + `). The cluster version is the upper
bound for every node pool's version, so pick it deliberately.

Public vs private:
  * Public (default): the Kubernetes API server is reachable on a public
    endpoint. You may still restrict access with --api-subnets.
  * Private (--public=false): the control plane stays on your network. This
    requires --location and --nat-gateway-ip (a reserved IP in that location),
    and optionally --node-subnet. These private-only properties are immutable.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the cluster reaches the AVAILABLE state.`,
		Example: `# Create a public cluster with the default Kubernetes version, then wait for it to become AVAILABLE
ionosctl compute k8s cluster create --name my-cluster --wait

# Create a public cluster on a specific version, restrict API-server access to two CIDRs, and ship audit logs to an S3 bucket
ionosctl compute k8s cluster create --name prod --k8s-version 1.29.5 --api-subnets 203.0.113.0/24,198.51.100.7/32 --s3bucket my-k8s-audit-logs

# Create a private cluster (control plane on your network) with a reserved NAT gateway IP
ionosctl compute k8s cluster create --name private --public=false --location de/txl --nat-gateway-ip 203.0.113.10 --node-subnet 10.0.0.0/16`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunK8sClusterCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "UnnamedCluster", "Name of the cluster. Up to 63 characters; must start and end with an alphanumeric character, dashes, underscores, dots and alphanumerics in between")
	cmd.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "Kubernetes control-plane version, e.g. 1.29.5. This is the upper bound for every node pool's version. If not set, the current default is used (see `ionosctl compute k8s version get`)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sVersion, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sVersionsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgS3Bucket, "", "", "Name of an existing IONOS Object Storage (S3) bucket that will receive Kubernetes API-server audit logs for this cluster")
	cmd.AddStringSliceFlag(cloudapiv6.ArgApiSubnets, "", []string{""}, "Restrict access to the Kubernetes API server to these CIDRs (allow list). Cluster-internal traffic is not affected. If left empty, access is not restricted. An IP given without a subnet mask defaults to /32 for IPv4 and /128 for IPv6. e.g. --api-subnets 203.0.113.0/24,198.51.100.7")
	cmd.AddBoolFlag(cloudapiv6.ArgPublic, "", true, "Whether the cluster is public (API server on a public endpoint, default) or private (--public=false, control plane on your network; requires --location and --nat-gateway-ip)")
	cmd.AddStringFlag(cloudapiv6.ArgLocation, "", "us/las", "Location of a private cluster (mandatory when --public=false, ignored for public clusters). Must be enabled for your contract or host a Data Center you own. Immutable. Location de/fra/2 is currently unavailable")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgNatGatewayIp, "", "", "NAT gateway IP for a private cluster (mandatory when --public=false, ignored for public clusters). Must be a reserved IP in the same location as --location. Immutable")
	cmd.AddStringFlag(constants.FlagNodeSubnet, "", "", "Node subnet for a private cluster (optional, ignored for public clusters). Must be a valid IPv4 CIDR with a /16 prefix, e.g. 10.0.0.0/16. Immutable")

	return cmd
}
