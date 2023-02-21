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

var tokenPutProperties = sdkgo.PostTokenProperties{}

func TokenPutCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "replace",
			Aliases:    []string{"r", "re"},
			ShortDesc:  "Create or replace a token",
			LongDesc:   "Create or replace a token used to access a container registry",
			Example:    "ionosctl container-registry token replace --name [NAME] --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]",
			PreCmdRun:  PreCmdPutToken,
			CmdRun:     CmdPutToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag("name", "", "", "Name of the Token", core.RequiredFlagOption())

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
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"enabled", "disabled",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdPutToken(c *core.CommandConfig) error {
	var err error

	regId, err := c.Command.Command.Flags().GetString("registry-id")
	if err != nil {
		return err
	}
	tokenId, err := c.Command.Command.Flags().GetString("token-id")
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString("name")
	if err != nil {
		return err
	}

	tokenPutProperties.SetName(name)

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
		tokenPutProperties.SetExpiryDate(expiryDate)

	}

	if viper.IsSet(core.GetFlagName(c.NS, "status")) {
		var status string
		status, err = c.Command.Command.Flags().GetString("status")
		if err != nil {
			return err
		}
		tokenPutProperties.SetStatus(status)

	}

	tokenInputPut := sdkgo.NewPutTokenInputWithDefaults()
	tokenInputPut.SetProperties(tokenPutProperties)

	token, _, err := c.ContainerRegistryServices.Token().Put(tokenId, *tokenInputPut, regId)
	if err != nil {
		return err
	}

	tokenPrint := sdkgo.NewTokenResponseWithDefaults()
	tokenPrint.SetProperties(*token.GetProperties())

	return c.Printer.Print(getTokenPrint(nil, c, &[]sdkgo.TokenResponse{*tokenPrint}, true))
}

func PreCmdPutToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", "registry-id", "name")
	if err != nil {
		return err
	}

	return nil
}
