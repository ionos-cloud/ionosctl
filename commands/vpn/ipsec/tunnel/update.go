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
			if c.Command.Command.Flags().Changed(constants.FlagJsonProperties) {
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
	gatewayID, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
	if err != nil {
		return err
	}
	tunnelID, err := c.Command.Command.Flags().GetString(constants.FlagTunnelID)
	if err != nil {
		return err
	}
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPut(context.Background(), gatewayID, tunnelID).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{Properties: propertiesFromJson}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func putFromProperties(c *core.CommandConfig) error {
	gatewayID, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
	if err != nil {
		return err
	}
	tunnelID, err := c.Command.Command.Flags().GetString(constants.FlagTunnelID)
	if err != nil {
		return err
	}
	original, _, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(context.Background(), gatewayID, tunnelID).
		Execute()
	if err != nil {
		return err
	}
	input := original.Properties

	if c.Command.Command.Flags().Changed(constants.FlagName) {
		name, err := c.Command.Command.Flags().GetString(constants.FlagName)
		if err != nil {
			return err
		}
		input.Name = name
	}

	if c.Command.Command.Flags().Changed(constants.FlagDescription) {
		desc, err := c.Command.Command.Flags().GetString(constants.FlagDescription)
		if err != nil {
			return err
		}
		input.Description = pointer.From(desc)
	}

	if c.Command.Command.Flags().Changed(constants.FlagHost) {
		host, err := c.Command.Command.Flags().GetString(constants.FlagHost)
		if err != nil {
			return err
		}
		input.RemoteHost = host
	}

	if c.Command.Command.Flags().Changed(constants.FlagAuthMethod) {
		method, err := c.Command.Command.Flags().GetString(constants.FlagAuthMethod)
		if err != nil {
			return err
		}
		input.Auth.Method = method
	}

	if c.Command.Command.Flags().Changed(constants.FlagPSKKey) {
		key, err := c.Command.Command.Flags().GetString(constants.FlagPSKKey)
		if err != nil {
			return err
		}
		input.Auth.Psk = &vpn.IPSecPSK{}
		input.Auth.Psk.Key = key
	}

	if c.Command.Command.Flags().Changed(constants.FlagIKEDiffieHellmanGroup) {
		dhg, err := c.Command.Command.Flags().GetString(constants.FlagIKEDiffieHellmanGroup)
		if err != nil {
			return err
		}
		input.Ike.DiffieHellmanGroup = pointer.From(dhg)
	}
	if c.Command.Command.Flags().Changed(constants.FlagIKEEncryptionAlgorithm) {
		alg, err := c.Command.Command.Flags().GetString(constants.FlagIKEEncryptionAlgorithm)
		if err != nil {
			return err
		}
		input.Ike.EncryptionAlgorithm = pointer.From(alg)
	}
	if c.Command.Command.Flags().Changed(constants.FlagIKEIntegrityAlgorithm) {
		alg, err := c.Command.Command.Flags().GetString(constants.FlagIKEIntegrityAlgorithm)
		if err != nil {
			return err
		}
		input.Ike.IntegrityAlgorithm = pointer.From(alg)
	}
	if c.Command.Command.Flags().Changed(constants.FlagIKELifetime) {
		lifetime, err := c.Command.Command.Flags().GetInt32(constants.FlagIKELifetime)
		if err != nil {
			return err
		}
		input.Ike.Lifetime = pointer.From(lifetime)
	}

	if c.Command.Command.Flags().Changed(constants.FlagESPDiffieHellmanGroup) {
		dhg, err := c.Command.Command.Flags().GetString(constants.FlagESPDiffieHellmanGroup)
		if err != nil {
			return err
		}
		input.Esp.DiffieHellmanGroup = pointer.From(dhg)
	}
	if c.Command.Command.Flags().Changed(constants.FlagESPEncryptionAlgorithm) {
		alg, err := c.Command.Command.Flags().GetString(constants.FlagESPEncryptionAlgorithm)
		if err != nil {
			return err
		}
		input.Esp.EncryptionAlgorithm = pointer.From(alg)
	}
	if c.Command.Command.Flags().Changed(constants.FlagESPIntegrityAlgorithm) {
		alg, err := c.Command.Command.Flags().GetString(constants.FlagESPIntegrityAlgorithm)
		if err != nil {
			return err
		}
		input.Esp.IntegrityAlgorithm = pointer.From(alg)
	}
	if c.Command.Command.Flags().Changed(constants.FlagESPLifetime) {
		lifetime, err := c.Command.Command.Flags().GetInt32(constants.FlagESPLifetime)
		if err != nil {
			return err
		}
		input.Esp.Lifetime = pointer.From(lifetime)
	}

	if c.Command.Command.Flags().Changed(constants.FlagCloudNetworkCIDRs) {
		cidrs, err := c.Command.Command.Flags().GetStringSlice(constants.FlagCloudNetworkCIDRs)
		if err != nil {
			return err
		}
		input.CloudNetworkCIDRs = cidrs
	}
	if c.Command.Command.Flags().Changed(constants.FlagPeerNetworkCIDRs) {
		cidrs, err := c.Command.Command.Flags().GetStringSlice(constants.FlagPeerNetworkCIDRs)
		if err != nil {
			return err
		}
		input.PeerNetworkCIDRs = cidrs
	}
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPut(context.Background(), gatewayID, tunnelID).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{
			Id:         tunnelID,
			Properties: input,
		}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}
