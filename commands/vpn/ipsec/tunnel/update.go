package tunnel

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
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
		ShortDesc: "Update a IPSec Tunnel",
		Example:   "ionosctl vpn ipsec tunnel update " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagTunnelID, constants.FlagName),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
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

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the IPSec Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL),
	)
	cmd.AddStringFlag(constants.FlagTunnelID, constants.FlagIdShort, "", "The ID of the IPSec Tunnel",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			gatewayID := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.TunnelIDs(gatewayID)
		}, constants.VPNApiRegionalURL),
	)

	cmd.AddStringFlag(constants.FlagName, "", "", "Name of the IPSec Tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Tunnel")
	cmd.AddStringFlag(constants.FlagHost, "", "", "The remote peer host fully qualified domain name or IPV4 IP to connect to. * __Note__: This should be the public IP of the remote peer. * Tunnels only support IPV4 or hostname (fully qualified DNS name).", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagAuthMethod, "", "", "The authentication method for the IPSec tunnel. Valid values are PSK or RSA", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagPSKKey, "", "", "The pre-shared key for the IPSec tunnel", core.RequiredFlagOption())

	cmd.AddSetFlag(constants.FlagIKEDiffieHellmanGroup, "", "", []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}, "The Diffie-Hellman Group to use for IPSec Encryption.")
	cmd.AddSetFlag(constants.FlagIKEEncryptionAlgorithm, "", "", []string{"AES128", "AES256"}, "The encryption algorithm to use for IPSec Encryption.")
	cmd.AddSetFlag(constants.FlagIKEIntegrityAlgorithm, "", "", []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}, "The integrity algorithm to use for IPSec Encryption.")
	cmd.AddInt32Flag(constants.FlagIKELifetime, "", 0, "The phase lifetime in seconds")

	cmd.AddSetFlag(constants.FlagESPDiffieHellmanGroup, "", "", []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}, "The Diffie-Hellman Group to use for IPSec Encryption.")
	cmd.AddSetFlag(constants.FlagESPEncryptionAlgorithm, "", "", []string{"AES128-CTR", "AES256-CTR", "AES128-GCM-16", "AES256-GCM-16", "AES128-GCM-12", "AES256-GCM-12", "AES128-CCM-12", "AES256-CCM-12", "AES128", "AES256"}, "The encryption algorithm to use for IPSec Encryption.")
	cmd.AddSetFlag(constants.FlagESPIntegrityAlgorithm, "", "", []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}, "The integrity algorithm to use for IPSec Encryption.")
	cmd.AddInt32Flag(constants.FlagESPLifetime, "", 0, "The phase lifetime in seconds")

	cmd.AddStringSliceFlag(constants.FlagCloudNetworkCIDRs, "", []string{}, "The network CIDRs on the \"Left\" side that are allowed to connect to the IPSec tunnel, i.e the CIDRs within your IONOS Cloud LAN. Specify \"0.0.0.0/0\" or \"::/0\" for all addresses.")
	cmd.AddStringSliceFlag(constants.FlagPeerNetworkCIDRs, "", []string{}, "The network CIDRs on the \"Right\" side that are allowed to connect to the IPSec tunnel. Specify \"0.0.0.0/0\" or \"::/0\" for all addresses.")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func putFromJSON(c *core.CommandConfig, propertiesFromJson vpn.IPSecTunnel) error {
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPut(context.Background(),
			viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{Properties: &propertiesFromJson}).Execute()
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
		input.Name = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
		input.Description = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagHost); viper.IsSet(fn) {
		input.RemoteHost = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagAuthMethod); viper.IsSet(fn) {
		input.Auth = &vpn.IPSecTunnelAuth{}
		input.Auth.Method = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagPSKKey); viper.IsSet(fn) {
		input.Auth.Psk = &vpn.IPSecPSK{}
		input.Auth.Psk.Key = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagIKEDiffieHellmanGroup); viper.IsSet(fn) {
		if input.Ike == nil {
			input.Ike = &vpn.IKEEncryption{}
		}
		input.Ike.DiffieHellmanGroup = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagIKEEncryptionAlgorithm); viper.IsSet(fn) {
		if input.Ike == nil {
			input.Ike = &vpn.IKEEncryption{}
		}
		input.Ike.EncryptionAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagIKEIntegrityAlgorithm); viper.IsSet(fn) {
		if input.Ike == nil {
			input.Ike = &vpn.IKEEncryption{}
		}
		input.Ike.IntegrityAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagIKELifetime); viper.IsSet(fn) {
		if input.Ike == nil {
			input.Ike = &vpn.IKEEncryption{}
		}
		input.Ike.Lifetime = pointer.From(int32(viper.GetInt(fn)))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagESPDiffieHellmanGroup); viper.IsSet(fn) {
		if input.Esp == nil {
			input.Esp = &vpn.ESPEncryption{}
		}
		input.Esp.DiffieHellmanGroup = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagESPEncryptionAlgorithm); viper.IsSet(fn) {
		if input.Esp == nil {
			input.Esp = &vpn.ESPEncryption{}
		}
		input.Esp.EncryptionAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagESPIntegrityAlgorithm); viper.IsSet(fn) {
		if input.Esp == nil {
			input.Esp = &vpn.ESPEncryption{}
		}
		input.Esp.IntegrityAlgorithm = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagESPLifetime); viper.IsSet(fn) {
		if input.Esp == nil {
			input.Esp = &vpn.ESPEncryption{}
		}
		input.Esp.Lifetime = pointer.From(int32(viper.GetInt(fn)))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagCloudNetworkCIDRs); viper.IsSet(fn) {
		input.CloudNetworkCIDRs = pointer.From(viper.GetStringSlice(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPeerNetworkCIDRs); viper.IsSet(fn) {
		input.PeerNetworkCIDRs = pointer.From(viper.GetStringSlice(fn))
	}
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPut(context.Background(),
			viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{
			Id:         pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))),
			Properties: input,
		}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}
