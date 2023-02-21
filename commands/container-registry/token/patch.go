package token

import (
	"context"
	"time"

	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenInput = sdkgo.NewPatchTokenInput()

func TokenPatchCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "update",
			Aliases:    []string{"u", "up"},
			ShortDesc:  "Update a token's properties",
			LongDesc:   "Use this command to update a token's properties. You can update the token's expiry date and status.",
			Example:    "ionosctl container-registry token update --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --expiry-date [EXPIRY-DATE] --status [STATUS]",
			PreCmdRun:  PreCmdPatchToken,
			CmdRun:     CmdPatchToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag("token-id", "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"token-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag("expiry-date", "", "", "Expiry date of the Token")
	cmd.AddStringFlag("status", "", "", "Status of the Token")

	return cmd
}

func CmdPatchToken(c *core.CommandConfig) error {
	var err error

	regId, err := c.Command.Command.Flags().GetString("registry-id")
	if err != nil {
		return err
	}
	tokenId, err := c.Command.Command.Flags().GetString("token-id")
	if err != nil {
		return err
	}

	if viper.IsSet(core.GetFlagName(c.NS, "expiry-date")) {
		var expiryDate time.Time
		expiryDateString, err := c.Command.Command.Flags().GetString("expiry-date")
		if err != nil {
			return err
		}
		expiryDate, err = time.Parse(time.RFC3339, expiryDateString)
		if err != nil {
			return err
		}
		tokenInput.SetExpiryDate(expiryDate)

	}

	if viper.IsSet(core.GetFlagName(c.NS, "status")) {
		var status string
		status, err = c.Command.Command.Flags().GetString("status")
		if err != nil {
			return err
		}
		tokenInput.SetStatus(status)

	}

	token, _, err := c.ContainerRegistryServices.Token().Patch(tokenId, *tokenInput, regId)
	if err != nil {
		return err
	}

	return c.Printer.Print(getTokenPrint(nil, c, &[]sdkgo.TokenResponse{token}, true))
}

func PreCmdPatchToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", "registry-id")
	if err != nil {
		return err
	}

	return nil
}
