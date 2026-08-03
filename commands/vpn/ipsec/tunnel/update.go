package tunnel

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/viper"
)

func Update() *core.Command {
	jsonPropertiesExample := "{\n  \"metadata\": {},\n  \"properties\": {\n    \"name\": \"My Company Gateway Tunnel\",\n    \"description\": \"Allows local subnet X to connect to virtual network Y.\",\n    \"remoteHost\": \"vpn.mycompany.com\",\n    \"auth\": {\n      \"method\": \"PSK\",\n      \"psk\": {\n        \"key\": \"X2wosbaw74M8hQGbK3jCCaEusR6CCFRa\"\n      }\n    },\n    \"ike\": {\n      \"diffieHellmanGroup\": \"16-MODP4096\",\n      \"encryptionAlgorithm\": \"AES256\",\n      \"integrityAlgorithm\": \"SHA256\",\n      \"lifetime\": 86400\n    },\n    \"esp\": {\n      \"diffieHellmanGroup\": \"16-MODP4096\",\n      \"encryptionAlgorithm\": \"AES256\",\n      \"integrityAlgorithm\": \"SHA256\",\n      \"lifetime\": 3600\n    },\n    \"cloudNetworkCIDRs\": [\n      \"192.168.1.100/24\"\n    ],\n    \"peerNetworkCIDRs\": [\n      \"1.2.3.4/32\"\n    ]\n  }\n}"
	tunnelViaJson := vpn.IPSecTunnel{}
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, jsonPropertiesExample, tunnelViaJson, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec tunnel",
		Verb:      "update",
		Aliases:   []string{"u", "patch", "put"},
		ShortDesc: "Update an IPSec Tunnel",
		LongDesc:  "Update an IPSec Tunnel. Any crypto or CIDR change must be mirrored on the remote peer or the tunnel will drop. See field meanings under 'ipsec tunnel create'.",
		Example:   "ionosctl vpn ipsec tunnel update " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagTunnelID, constants.FlagName),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsSetsAndLocation(
				[]string{constants.FlagGatewayID, constants.FlagTunnelID},
				[]string{constants.FlagJsonProperties, constants.FlagGatewayID, constants.FlagTunnelID},
				[]string{constants.FlagJsonPropertiesExample},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.IsSet(constants.FlagJsonProperties) {
				return putFromJSON(c, tunnelViaJson)
			}

			return putFromProperties(c)
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the IPSec Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)
	cmd.AddStringFlag(constants.FlagTunnelID, constants.FlagIdShort, "", "The ID of the IPSec Tunnel",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			gatewayID := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.TunnelIDs(gatewayID)
		}, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the IPSec Tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Tunnel")
	cmd.AddStringFlag(constants.FlagHost, "", "", "Public IPv4 or fully-qualified hostname of the remote peer to connect to (the remote side's public address; IPv6 is not supported)", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagAuthMethod, "", "", "How the two ends authenticate each other: PSK (shared secret in --psk-key) or RSA", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagPSKKey, "", "", "Pre-shared key, when --auth-method is PSK; the identical secret must be configured on the remote peer", core.RequiredFlagOption())

	cmd.AddSetFlag(constants.FlagIKEDiffieHellmanGroup, "", "", []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}, "IKE (phase 1) Diffie-Hellman group for the key exchange; must match the remote peer")
	cmd.AddSetFlag(constants.FlagIKEEncryptionAlgorithm, "", "", []string{"AES128", "AES256"}, "IKE (phase 1) encryption algorithm; must match the remote peer")
	cmd.AddSetFlag(constants.FlagIKEIntegrityAlgorithm, "", "", []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}, "IKE (phase 1) integrity/hash algorithm; must match the remote peer")
	cmd.AddInt32Flag(constants.FlagIKELifetime, "", 0, "IKE (phase 1) rekey interval in seconds; 0 uses the API default (e.g. 86400 = 24h)")

	cmd.AddSetFlag(constants.FlagESPDiffieHellmanGroup, "", "", []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}, "ESP (phase 2) Diffie-Hellman group for the data channel; must match the remote peer")
	cmd.AddSetFlag(constants.FlagESPEncryptionAlgorithm, "", "", []string{"AES128-CTR", "AES256-CTR", "AES128-GCM-16", "AES256-GCM-16", "AES128-GCM-12", "AES256-GCM-12", "AES128-CCM-12", "AES256-CCM-12", "AES128", "AES256"}, "ESP (phase 2) encryption algorithm for the data channel; must match the remote peer")
	cmd.AddSetFlag(constants.FlagESPIntegrityAlgorithm, "", "", []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}, "ESP (phase 2) integrity/hash algorithm; must match the remote peer")
	cmd.AddInt32Flag(constants.FlagESPLifetime, "", 0, "ESP (phase 2) rekey interval in seconds; 0 uses the API default (e.g. 3600 = 1h)")

	cmd.AddStringSliceFlag(constants.FlagCloudNetworkCIDRs, "", []string{}, "Local IONOS-side subnets (CIDR) allowed to cross the tunnel, i.e. the networks in your IONOS CLOUD LAN. Use \"0.0.0.0/0\",\"::/0\" for all addresses")
	cmd.AddStringSliceFlag(constants.FlagPeerNetworkCIDRs, "", []string{}, "Remote-side subnets (CIDR) reachable through the tunnel, i.e. the networks behind the remote peer. Use \"0.0.0.0/0\",\"::/0\" for all addresses")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func putFromJSON(c *core.CommandConfig, propertiesFromJson vpn.IPSecTunnel) error {
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPut(context.Background(),
			viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{Properties: propertiesFromJson}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func putFromProperties(c *core.CommandConfig) error {
	original, _, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))).
		Execute()
	if err != nil {
		return err
	}
	input := original.Properties

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
		input.Auth.Method = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, constants.FlagPSKKey); viper.IsSet(fn) {
		input.Auth.Psk = &vpn.IPSecPSK{}
		input.Auth.Psk.Key = viper.GetString(fn)
	}

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
		IpsecgatewaysTunnelsPut(context.Background(),
			viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{
			Id:         viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID)),
			Properties: input,
		}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}
