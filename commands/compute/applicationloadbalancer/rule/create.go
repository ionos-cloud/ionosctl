package rule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApplicationLoadBalancerForwardingRuleCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to create a forwarding rule (a listener) on a specified Application Load Balancer. The rule binds a protocol, an inbound IP and a port that the balancer will accept client connections on.

The --listener-ip must be one of the ALB's own IPs on its listener LAN. For an HTTPS listener, additionally pass --server-certificates so the balancer can present a certificate during the TLS handshake. After the rule exists, attach HTTP rules to it (` + "`" + `alb rule httprule add` + "`" + `) to define how matching requests are routed.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Listener Ip
* Listener Port`,
		Example: `# Create an HTTP listener on port 80
ionosctl compute applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --name "http-listener" --listener-ip 192.0.2.10 --listener-port 80

# Create an HTTPS listener on port 443 with a server certificate and a longer client timeout
ionosctl compute applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --name "https-listener" --listener-ip 192.0.2.10 --listener-port 443 --server-certificates CERTIFICATE_ID --client-timeout 60000 --wait`,
		PreCmdRun:  PreRunApplicationLoadBalancerForwardingRuleCreate,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.Name(), cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Forwarding Rule", "The name of the Application Load Balancer forwarding rule.")
	cmd.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "The listener protocol. HTTP is the only supported value (the ALB is a layer-7 balancer); HTTPS listeners are configured by using HTTP together with --server-certificates.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "The inbound IP the balancer listens on. Must be one of the ALB's own --ips assigned on its listener LAN.", core.RequiredFlagOption())
	cmd.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "The inbound TCP port the balancer listens on; valid range is 1 to 65535 (typically 80 for HTTP or 443 for HTTPS).", core.RequiredFlagOption())
	cmd.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 50, "The maximum time in milliseconds to wait for the client to acknowledge or send data before the connection is closed.")
	cmd.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "IDs of server certificates (managed by the IONOS Certificate Manager) that the balancer presents to clients during the TLS handshake. Required to serve HTTPS on this listener.")
	cmd.AddColsFlag(allAlbForwardingRuleCols)

	return cmd
}
