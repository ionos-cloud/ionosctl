package token

import (
	"context"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenInput = containerregistry.NewPatchTokenInput()

func TokenUpdateCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "token",
			Verb:      "update",
			Aliases:   []string{"u", "up"},
			ShortDesc: "Update a token's expiry or status (PATCH)",
			LongDesc: `Update an existing token in place (HTTP PATCH). Only the fields you pass are changed; the token's scopes and password are preserved (unlike 'replace', which regenerates the password and clears scopes).

Use --status disabled to revoke a token without deleting it (and enabled to re-activate it), or change its expiry with --expiry-date (absolute RFC3339) / --expiry-time (relative duration). To change what the token may access, use 'container-registry token scope'.`,
			Example: `# Revoke a token without deleting it
ionosctl container-registry token update --registry-id REGISTRY_ID --token-id TOKEN_ID --status disabled

# Extend a token's expiry to an absolute date
ionosctl container-registry token update --registry-id REGISTRY_ID --token-id TOKEN_ID --expiry-date 2026-01-01T00:00:00Z`,
			PreCmdRun:  PreCmdPatchToken,
			CmdRun:     CmdPatchToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry that owns the token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "", "", "The unique ID of the token to update")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId,
		func(cobracmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(
				viper.GetString(
					core.GetFlagName(
						cmd.NS, constants.FlagRegistryId,
					),
				),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(FlagExpiryDate, "", "", "New absolute expiry date as an RFC3339 timestamp, e.g. 2025-01-02T15:04:05Z")
	cmd.AddStringFlag(FlagTimeUntilExpiry, "", "", "New expiry as a duration from now, combining y/m/d/h, e.g. 1y2d")
	cmd.AddStringFlag(FlagStatus, "", "", "Token status: 'enabled' (usable) or 'disabled' (revoked)")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"enabled", "disabled",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdPatchToken(c *core.CommandConfig) error {
	var err error

	regId, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
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

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensPatch(context.Background(), regId, tokenId).PatchTokenInput(*tokenInput).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(token)
}

func PreCmdPatchToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagTokenId, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
