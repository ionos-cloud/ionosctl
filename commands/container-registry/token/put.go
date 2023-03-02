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

	cmd.AddStringFlag(FlagName, "", "", "Name of the Token", core.RequiredFlagOption())

	cmd.AddStringFlag(FlagRegId, "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag("token-id", "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"token-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(FlagExpiryDate, "", "", "Expiry date of the Token")
	cmd.AddStringFlag(FlagTimeUntilExpiry, "", "", "Time until the Token expires (ex: 1y2d)")
	cmd.AddStringFlag(FlagStatus, "", "", "Status of the Token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"enabled", "disabled",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.MarkFlagsMutuallyExclusive(FlagExpiryDate, FlagTimeUntilExpiry)
	return cmd
}

func CmdPutToken(c *core.CommandConfig) error {
	var err error

	regId, err := c.Command.Command.Flags().GetString(FlagRegId)
	if err != nil {
		return err
	}
	tokenId, err := c.Command.Command.Flags().GetString("token-id")
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(FlagName)
	if err != nil {
		return err
	}

	tokenPutProperties.SetName(name)

	if viper.IsSet(core.GetFlagName(c.NS, FlagExpiryDate)) {
		var expiryDate time.Time
		expiryDateString, err := c.Command.Command.Flags().GetString(FlagExpiryDate)
		if err != nil {
			return err
		}
		expiryDate, err = time.Parse(time.RFC3339, expiryDateString)
		if err != nil {
			return err
		}
		tokenPutProperties.SetExpiryDate(expiryDate)

	} else if viper.IsSet(core.GetFlagName(c.NS, FlagTimeUntilExpiry)) {
		var timeUntilExpiry string
		timeUntilExpiry, err = c.Command.Command.Flags().GetString(FlagTimeUntilExpiry)
		if err != nil {
			return err
		}
		timeNow := time.Now()
		duration, err := ParseExpiryTime(timeUntilExpiry)
		if err != nil {
			return err
		}
		timeNow.Add(duration)
	}

	if viper.IsSet(core.GetFlagName(c.NS, FlagStatus)) {
		var status string
		status, err = c.Command.Command.Flags().GetString(FlagStatus)
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
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", FlagRegId, FlagName)
	if err != nil {
		return err
	}

	return nil
}
