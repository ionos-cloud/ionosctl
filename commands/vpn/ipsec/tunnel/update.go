package tunnel

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
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

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the IPSec Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GatewayIDs(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagTunnelID, constants.FlagIdShort, "", "The ID of the IPSec Tunnel you want to delete", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagTunnelID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return TunnelsProperty(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagTunnelID)), func(p vpn.IPSecTunnelRead) string {
			return *p.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, "", "", "Name of the IPSec Tunnel")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Tunnel")
	cmd.AddStringSliceFlag(constants.FlagIps, "", []string{}, "Comma separated subnets of CIDRs that are allowed to connect to the IPSec Gateway. Specify \"a.b.c.d/32\" for an individual IP address. Specify \"0.0.0.0/0\" or \"::/0\" for all addresses")
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagIps, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"::/0"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagPublicKey, "", "", "Public key of the connecting tunnel")
	cmd.AddStringFlag(constants.FlagHost, "", "", "Hostname or IPV4 address that the IPSec Server will connect to")
	cmd.AddIntFlag(constants.FlagPort, "", 51820, "Port that the IPSec Server will connect to")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func putFromJSON(c *core.CommandConfig, propertiesFromJson vpn.IPSecTunnel) error {
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPut(context.Background(),
			viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(constants.FlagTunnelID)).
		IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{Properties: &propertiesFromJson}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func putFromProperties(c *core.CommandConfig) error {
	original, _, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)), viper.GetString(constants.FlagTunnelID)).
		Execute()
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

	input.Ike = &vpn.IKEEncryption{}
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

	input.Esp = &vpn.ESPEncryption{}
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
