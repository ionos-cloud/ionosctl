package token

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenPutProperties = sdkgo.PostTokenProperties{}

func TokenReplaceCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "token",
			Verb:      "replace",
			Aliases:   []string{"r", "re"},
			ShortDesc: "Create or replace a token",
			LongDesc:  "Create or replace a token used to access a container registry",
			Example: "ionosctl container-registry token replace --name [NAME] --registry-id [REGISTRY-ID] --token-id [TOKEN-ID]\n" +
				"In order to save the token to a environment variable: export [ENV-VAL-NAME]=$(ionosctl cr token replace --name [TOKEN-NAME] --registry-id [REGISTRY-ID] --token-id [TOKEN-ID]",
			PreCmdRun:  PreCmdPutToken,
			CmdRun:     CmdPutToken,
			InitClient: true,
		},
	)

	// This line is only used to override the help text for `--no-headers`!
	cmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, true, "Use --no-headers=false to show column headers",
	)

	cmd.AddStringFlag(FlagName, "", "", "Name of the Token", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagRegId, "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId, func(cobracmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(viper.GetString(core.GetFlagName(cmd.NS, FlagRegId))), cobra.ShellCompDirectiveNoFileComp
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
	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(AllTokenCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return AllTokenCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.MarkFlagsMutuallyExclusive(FlagExpiryDate, FlagTimeUntilExpiry)
	return cmd
}

func CmdPutToken(c *core.CommandConfig) error {
	if !viper.IsSet(constants.ArgNoHeaders) {
		viper.Set(constants.ArgNoHeaders, true) // Change default to work as for `token create`
	}

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

		timeNow = timeNow.Add(duration)
		tokenPutProperties.SetExpiryDate(timeNow)
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

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ContainerRegistryToken, token, tabheaders.GetHeaders(AllTokenCols, postHeaders, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func PreCmdPutToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", FlagRegId, FlagName)
	if err != nil {
		return err
	}

	return nil
}
