package token

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenInput = sdkgo.NewPatchTokenInput()

func TokenUpdateCmd() *core.Command {
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
	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(AllTokenCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return AllTokenCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdPatchToken(c *core.CommandConfig) error {
	var err error

	regId, err := c.Command.Command.Flags().GetString(FlagRegId)
	if err != nil {
		return err
	}

	tokenId, err := c.Command.Command.Flags().GetString(FlagTokenId)
	if err != nil {
		return err
	}

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

		tokenInput.SetExpiryDate(expiryDate)
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

		tokenInput.SetStatus(status)
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	token, _, err := c.ContainerRegistryServices.Token().Patch(tokenId, *tokenInput, regId)
	if err != nil {
		return err
	}
	out, err := jsontabwriter.GenerateOutput("", allJSONPaths, token, tabheaders.GetHeadersAllDefault(AllTokenCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func PreCmdPatchToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagTokenId, FlagRegId)
	if err != nil {
		return err
	}

	return nil
}
