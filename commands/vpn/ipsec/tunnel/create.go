package tunnel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	jsonPropertiesExample := "{\n  \"properties\": {\n    \"name\": \"My Company Gateway Tunnel\",\n    \"description\": \"Allows local subnet X to connect to virtual network Y.\",\n    \"remoteHost\": \"vpn.mycompany.com\",\n    \"auth\": {\n      \"method\": \"PSK\",\n      \"psk\": {\n        \"key\": \"X2wosbaw74M8hQGbK3jCCaEusR6CCFRa\"\n      }\n    },\n    \"ike\": {\n      \"diffieHellmanGroup\": \"16-MODP4096\",\n      \"encryptionAlgorithm\": \"AES256\",\n      \"integrityAlgorithm\": \"SHA256\",\n      \"lifetime\": 86400\n    },\n    \"esp\": {\n      \"diffieHellmanGroup\": \"16-MODP4096\",\n      \"encryptionAlgorithm\": \"AES256\",\n      \"integrityAlgorithm\": \"SHA256\",\n      \"lifetime\": 3600\n    },\n    \"cloudNetworkCIDRs\": [\n      \"192.168.1.100/24\"\n    ],\n    \"peerNetworkCIDRs\": [\n      \"1.2.3.4/32\"\n    ]\n  }\n}"
	tunnelViaJson := vpn.IPSecTunnelCreate{}
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, jsonPropertiesExample, &tunnelViaJson,
		core.CommandBuilder{
			Namespace: "vpn",
			Resource:  "ipsec tunnel",
			Verb:      "create",
			Aliases:   []string{"c", "post"},
			ShortDesc: "Create an IPSec tunnel",
			LongDesc: `Create an IPSec tunnel: the connection from a gateway (--gateway-id) to one remote site.

Point --host at the remote peer's public IPv4 or FQDN, then authenticate: --auth-method PSK together with a shared --psk-key (or RSA). Set the phase-1 (IKE) and phase-2 (ESP) crypto with the --ike-* / --esp-* flags — each takes a Diffie-Hellman group, an encryption algorithm, an integrity algorithm and a lifetime in seconds (rekey interval; leave 0 to use the API default). Finally list which subnets may cross: --cloud-network-cidrs on your IONOS LAN side and --peer-network-cidrs on the remote side.

Both ends must use the SAME crypto parameters and mirrored CIDRs or the tunnel stays down.

You can instead pass the whole request body with --json-properties (see --json-properties-example for a template).`,
			Example: `ionosctl vpn ipsec tunnel create --gateway-id GATEWAY_ID --name to-hq --host vpn.example.com --auth-method PSK --psk-key SHARED_SECRET --ike-diffie-hellman-group 16-MODP4096 --ike-encryption-algorithm AES256 --ike-integrity-algorithm SHA256 --esp-diffie-hellman-group 16-MODP4096 --esp-encryption-algorithm AES256 --esp-integrity-algorithm SHA256 --cloud-network-cidrs 10.7.222.0/24 --peer-network-cidrs 192.168.1.0/24
ionosctl vpn ipsec tunnel create --gateway-id GATEWAY_ID --json-properties tunnel.json
ionosctl vpn ipsec tunnel create --json-properties-example`,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return c.CheckRequiredFlagsSetsAndLocation(
					[]string{constants.FlagJsonProperties, constants.FlagGatewayID},
					[]string{constants.FlagJsonPropertiesExample},
					[]string{
						constants.FlagGatewayID,
						constants.FlagName,
						constants.FlagHost,
						constants.FlagAuthMethod,
						constants.FlagPSKKey,
						constants.FlagIKEDiffieHellmanGroup,
						constants.FlagIKEEncryptionAlgorithm,
						constants.FlagIKEIntegrityAlgorithm,
						constants.FlagIKELifetime,
						constants.FlagESPDiffieHellmanGroup,
						constants.FlagESPEncryptionAlgorithm,
						constants.FlagESPIntegrityAlgorithm,
						constants.FlagESPLifetime,
						constants.FlagCloudNetworkCIDRs,
						constants.FlagPeerNetworkCIDRs,
					},
				)
			},
			CmdRun: func(c *core.CommandConfig) error {
				if c.Command.Command.Flags().Changed(constants.FlagJsonProperties) {
					j, _ := json.MarshalIndent(tunnelViaJson, "", "    ")
					fmt.Println(string(j))

					return createFromJSON(c, tunnelViaJson)
				}
				return createFromProperties(c)
			},
		})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the IPSec Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the IPSec Tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Tunnel")
	cmd.AddStringFlag(constants.FlagHost, "", "", "Public IPv4 or fully-qualified hostname of the remote peer to connect to (the remote side's public address; IPv6 is not supported)", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagAuthMethod, "", "", "How the two ends authenticate each other: PSK (shared secret in --psk-key) or RSA", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagAuthMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"PSK", "RSA"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagPSKKey, "", "", "Pre-shared key, when --auth-method is PSK; the identical secret must be configured on the remote peer", core.RequiredFlagOption())

	cmd.AddSetFlag(constants.FlagIKEDiffieHellmanGroup, "", "", []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}, "IKE (phase 1) Diffie-Hellman group for the key exchange; must match the remote peer")
	cmd.AddSetFlag(constants.FlagIKEEncryptionAlgorithm, "", "", []string{"AES128", "AES256"}, "IKE (phase 1) encryption algorithm; must match the remote peer")
	cmd.AddSetFlag(constants.FlagIKEIntegrityAlgorithm, "", "", []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}, "IKE (phase 1) integrity/hash algorithm; must match the remote peer")
	cmd.AddInt32Flag(constants.FlagIKELifetime, "", 0, "IKE (phase 1) rekey interval in seconds; 0 uses the API default (e.g. 86400 = 24h)")

	cmd.AddSetFlag(constants.FlagESPDiffieHellmanGroup, "", "", []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}, "ESP (phase 2) Diffie-Hellman group for the data channel; must match the remote peer")
	cmd.AddSetFlag(constants.FlagESPEncryptionAlgorithm, "", "", []string{"AES128-CTR", "AES256-CTR", "AES128-GCM-16", "AES256-GCM-16", "AES128-GCM-12", "AES256-GCM-12", "AES128-CCM-12", "AES256-CCM-12", "AES128", "AES256"}, "ESP (phase 2) encryption algorithm for the data channel; must match the remote peer")
	cmd.AddSetFlag(constants.FlagESPIntegrityAlgorithm, "", "", []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}, "ESP (phase 2) integrity/hash algorithm; must match the remote peer")
	cmd.AddInt32Flag(constants.FlagESPLifetime, "", 0, "ESP (phase 2) rekey interval in seconds; 0 uses the API default (e.g. 3600 = 1h)")

	cmd.AddStringSliceFlag(constants.FlagCloudNetworkCIDRs, "", []string{}, "Local IONOS-side subnets (CIDR) allowed to cross the tunnel, i.e. the networks in your IONOS Cloud LAN. Use \"0.0.0.0/0\",\"::/0\" for all addresses")
	cmd.AddStringSliceFlag(constants.FlagPeerNetworkCIDRs, "", []string{}, "Remote-side subnets (CIDR) reachable through the tunnel, i.e. the networks behind the remote peer. Use \"0.0.0.0/0\",\"::/0\" for all addresses")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func createFromJSON(c *core.CommandConfig, propertiesFromJson vpn.IPSecTunnelCreate) error {
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPost(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))).
		IPSecTunnelCreate(propertiesFromJson).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func createFromProperties(c *core.CommandConfig) error {
	input := vpn.IPSecTunnel{}

	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
		input.Description = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagHost); viper.IsSet(fn) {
		input.RemoteHost = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, constants.FlagAuthMethod); viper.IsSet(fn) {
		input.Auth = vpn.IPSecTunnelAuth{}
		input.Auth.Method = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, constants.FlagPSKKey); viper.IsSet(fn) {
		input.Auth.Psk = &vpn.IPSecPSK{}
		input.Auth.Psk.Key = viper.GetString(fn)
	}

	input.Ike = vpn.IKEEncryption{}
	if fn := core.GetFlagName(c.NS, constants.FlagIKEDiffieHellmanGroup); viper.IsSet(fn) {
		input.Ike.DiffieHellmanGroup = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagIKEEncryptionAlgorithm); viper.IsSet(fn) {
		input.Ike.EncryptionAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagIKEIntegrityAlgorithm); viper.IsSet(fn) {
		input.Ike.IntegrityAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagIKELifetime); viper.IsSet(fn) {
		input.Ike.Lifetime = pointer.From(int32(viper.GetInt(fn)))
	}

	input.Esp = vpn.ESPEncryption{}
	if fn := core.GetFlagName(c.NS, constants.FlagESPDiffieHellmanGroup); viper.IsSet(fn) {
		input.Esp.DiffieHellmanGroup = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagESPEncryptionAlgorithm); viper.IsSet(fn) {
		input.Esp.EncryptionAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagESPIntegrityAlgorithm); viper.IsSet(fn) {
		input.Esp.IntegrityAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagESPLifetime); viper.IsSet(fn) {
		input.Esp.Lifetime = pointer.From(int32(viper.GetInt(fn)))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagCloudNetworkCIDRs); viper.IsSet(fn) {
		input.CloudNetworkCIDRs = viper.GetStringSlice(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPeerNetworkCIDRs); viper.IsSet(fn) {
		input.PeerNetworkCIDRs = viper.GetStringSlice(fn)
	}
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPost(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))).
		IPSecTunnelCreate(vpn.IPSecTunnelCreate{Properties: input}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func handleOutput(c *core.CommandConfig, tunnel vpn.IPSecTunnelRead) error {
	return c.Printer(allCols).Print(tunnel)
}
