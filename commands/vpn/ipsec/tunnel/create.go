package tunnel

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/gateway"
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
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec tunnel",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a IPSec tunnel",
		LongDesc:  "Create IPSec tunnels",
		Example:   "", // TODO: Probably best if I don't forget this
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagGatewayID,
				constants.FlagName,
				constants.FlagHost,
				constants.FlagAuthMethod,
				constants.FlagPSKKey,
				constants.FlagIKEEncryptionAlgorithm,
				constants.FlagIKEIntegrityAlgorithm,
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
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
			if fn := core.GetFlagName(c.NS, constants.FlagIKEEncryptionAlgorithm); viper.IsSet(fn) {
				input.Ike.EncryptionAlgorithm = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagIKEIntegrityAlgorithm); viper.IsSet(fn) {
				input.Ike.IntegrityAlgorithm = pointer.From(viper.GetString(fn))
			}

			tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
				IpsecgatewaysTunnelsPost(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))).
				IPSecTunnelCreate(vpn.IPSecTunnelCreate{Properties: input}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.VPNIPSecTunnel, tunnel, tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the IPSec Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return gateway.GatewaysProperty(func(gateway vpn.IPSecGatewayRead) string {
			return *gateway.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, "", "", "Name of the IPSec Tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Tunnel")
	cmd.AddStringFlag(constants.FlagHost, "", "", "The remote peer host fully qualified domain name or IPV4 IP to connect to. * __Note__: This should be the public IP of the remote peer. * Tunnels only support IPV4 or hostname (fully qualified DNS name).", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagAuthMethod, "", "", "The authentication method for the IPSec tunnel. Valid values are PSK or RSA", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagPSKKey, "", "", "The pre-shared key for the IPSec tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagIKEEncryptionAlgorithm, "", "", "The encryption algorithm for the IPSec tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagIKEIntegrityAlgorithm, "", "", "The integrity algorithm for the IPSec tunnel", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
