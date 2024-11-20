package tunnel

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	jsonPropertiesExample := "{\n  \"properties\": {\n    \"name\": \"My Company Gateway Tunnel\",\n    \"description\": \"Allows local subnet X to connect to virtual network Y.\",\n    \"remoteHost\": \"vpn.mycompany.com\",\n    \"auth\": {\n      \"method\": \"PSK\",\n      \"psk\": {\n        \"key\": \"X2wosbaw74M8hQGbK3jCCaEusR6CCFRa\"\n      }\n    },\n    \"ike\": {\n      \"diffieHellmanGroup\": \"16-MODP4096\",\n      \"encryptionAlgorithm\": \"AES256\",\n      \"integrityAlgorithm\": \"SHA256\",\n      \"lifetime\": 86400\n    },\n    \"esp\": {\n      \"diffieHellmanGroup\": \"16-MODP4096\",\n      \"encryptionAlgorithm\": \"AES256\",\n      \"integrityAlgorithm\": \"SHA256\",\n      \"lifetime\": 3600\n    },\n    \"cloudNetworkCIDRs\": [\n      \"192.168.1.100/24\"\n    ],\n    \"peerNetworkCIDRs\": [\n      \"1.2.3.4/32\"\n    ]\n  }\n}"
	tunnelViaJson := vpn.IPSecTunnel{}
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, jsonPropertiesExample, tunnelViaJson,
		core.CommandBuilder{
			Namespace: "vpn",
			Resource:  "ipsec tunnel",
			Verb:      "create",
			Aliases:   []string{"c", "post"},
			ShortDesc: "Create a IPSec tunnel",
			LongDesc:  "Create IPSec tunnels",
			Example:   "ionosctl vpn ipsec tunnel create " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagName, constants.FlagHost, constants.FlagAuthMethod, constants.FlagPSKKey, constants.FlagIKEDiffieHellmanGroup, constants.FlagIKEEncryptionAlgorithm, constants.FlagIKEIntegrityAlgorithm, constants.FlagIKELifetime, constants.FlagESPDiffieHellmanGroup, constants.FlagESPEncryptionAlgorithm, constants.FlagESPIntegrityAlgorithm, constants.FlagESPLifetime, constants.FlagCloudNetworkCIDRs, constants.FlagPeerNetworkCIDRs) + "\n" + "ionosctl vpn ipsec tunnel create " + core.FlagsUsage(constants.FlagJsonProperties) + "\n" + "ionosctl vpn ipsec tunnel create " + core.FlagsUsage(constants.FlagJsonProperties) + " " + constants.FlagJsonPropertiesExample,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlagsSets(c.Command, c.NS,
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
				if viper.IsSet(constants.FlagJsonProperties) {
					return createFromJSON(c, tunnelViaJson)
				}
				return createFromProperties(c)
			},
		})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the IPSec Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GatewayIDs(), cobra.ShellCompDirectiveNoFileComp
	})

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

func createFromJSON(c *core.CommandConfig, propertiesFromJson vpn.IPSecTunnel) error {
	tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
		IpsecgatewaysTunnelsPost(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))).
		IPSecTunnelCreate(vpn.IPSecTunnelCreate{Properties: &propertiesFromJson}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func createFromProperties(c *core.CommandConfig) error {
	input := &vpn.IPSecTunnel{}

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
		IpsecgatewaysTunnelsPost(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))).
		IPSecTunnelCreate(vpn.IPSecTunnelCreate{Properties: input}).Execute()
	if err != nil {
		return err
	}

	return handleOutput(c, tunnel)
}

func handleOutput(c *core.CommandConfig, tunnel vpn.IPSecTunnelRead) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.VPNIPSecTunnel, tunnel, tabheaders.GetHeaders(allCols, defaultCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
