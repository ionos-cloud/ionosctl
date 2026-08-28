package token

import (
	"context"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenPutProperties = containerregistry.PostTokenProperties{}

func TokenReplaceCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "token",
			Verb:      "replace",
			Aliases:   []string{"r", "re"},
			ShortDesc: "Replace a token, regenerating its password (PUT)",
			LongDesc: `Replace an existing token (HTTP PUT). This overwrites the whole token with the values you pass, so any omitted property (expiry, status, and its scopes) is reset to its default - in particular, existing scopes are cleared. To keep scopes, prefer 'container-registry token update' or re-add them with 'container-registry token scope add'.

Replacing regenerates the token password, which is printed only once in the response (capture it - see the second example). The old password stops working. Set expiry with --expiry-date (absolute RFC3339) or --expiry-time (relative duration); the two are mutually exclusive.`,
			Example: `# Replace a token (prints a new, one-time password)
ionosctl container-registry token replace --registry-id REGISTRY_ID --token-id TOKEN_ID --name push-token

# Replace and capture the new password into an env var
export CR_TOKEN=$(ionosctl cr token replace --registry-id REGISTRY_ID --token-id TOKEN_ID --name ci-token --expiry-time 1y)`,
			PreCmdRun:  PreCmdPutToken,
			CmdRun:     CmdPutToken,
			InitClient: true,
		},
	)

	// This line is only used to override the help text for `--no-headers`!
	cmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, true, "Use --no-headers=false to show column headers",
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the token", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry that owns the token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "", "", "The unique ID of the token to replace")
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

	cmd.AddStringFlag(FlagExpiryDate, "", "", "Absolute expiry date as an RFC3339 timestamp, e.g. 2025-01-02T15:04:05Z. Mutually exclusive with --expiry-time; omit both to never expire")
	cmd.AddStringFlag(FlagTimeUntilExpiry, "", "", "Relative expiry as a duration from now, combining y/m/d/h, e.g. 1y2d. Mutually exclusive with --expiry-date")
	cmd.AddStringFlag(FlagStatus, "", "", "Token status: 'enabled' (usable) or 'disabled' (revoked)")
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
	if !viper.IsSet(constants.ArgNoHeaders) {
		viper.Set(constants.ArgNoHeaders, true)
	}

	var err error

	regId, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	tokenId, err := c.Command.Command.Flags().GetString("token-id")
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
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

	tokenInputPut := containerregistry.NewPutTokenInputWithDefaults()
	tokenInputPut.SetProperties(tokenPutProperties)

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensPut(context.Background(), regId, tokenId).PutTokenInput(*tokenInputPut).Execute()
	if err != nil {
		return err
	}

	cols := c.Cols()
	// Default to showing only CredentialsPassword for replace operations
	if cols == nil {
		cols = []string{"CredentialsPassword"}
	}
	return c.Out(table.Sprint(allCols, token, cols))
}

func PreCmdPutToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", constants.FlagRegistryId, constants.FlagName)
	if err != nil {
		return err
	}

	return nil
}
