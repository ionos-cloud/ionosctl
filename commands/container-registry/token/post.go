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

var tokenPostProperties = containerregistry.NewPostTokenPropertiesWithDefaults()

func TokenPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "token",
			Verb:      "create",
			Aliases:   []string{"c"},
			ShortDesc: "Create a token to authenticate against a registry",
			LongDesc: `Create a new token used to authenticate 'docker login'/'docker push'/'docker pull' against a registry.

The generated password is printed only once, in this response, and cannot be retrieved afterwards - capture it now (see the second example). A freshly created token has no scopes yet, so it cannot pull or push until you grant scopes with 'container-registry token scope add'.

Set an expiry with either --expiry-date (an absolute RFC3339 timestamp) or --expiry-time (a relative duration such as 1y2d); the two are mutually exclusive. Omit both for a token that never expires. Use --status disabled to create the token pre-revoked.`,
			Example: `# Create a token (note the printed password - it is shown only once)
ionosctl container-registry token create --registry-id REGISTRY_ID --name push-token

# Create a token expiring in 1 year and 2 days, capturing the password into an env var
export CR_TOKEN=$(ionosctl cr token create --registry-id REGISTRY_ID --name ci-token --expiry-time 1y2d)`,
			PreCmdRun:  PreCmdPostToken,
			CmdRun:     CmdPostToken,
			InitClient: true,
		},
	)

	// This line is only used to override the help text for `--no-headers`!
	cmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, true, "Use --no-headers=false to show column headers",
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the token, unique within the registry", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagExpiryDate, "", "", "Absolute expiry date as an RFC3339 timestamp, e.g. 2025-01-02T15:04:05Z. Mutually exclusive with --expiry-time; omit both to never expire")
	cmd.AddStringFlag(FlagStatus, "", "", "Initial status of the token: 'enabled' (usable) or 'disabled' (revoked). Defaults to enabled")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"enabled", "disabled",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTimeUntilExpiry, "", "", "Relative expiry as a duration from now, combining y (years), m (months), d (days), h (hours), e.g. 1y2d or 6m. Mutually exclusive with --expiry-date")
	cmd.AddStringFlag(
		constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry that will own this token", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.MarkFlagsMutuallyExclusive(FlagExpiryDate, FlagTimeUntilExpiry)
	return cmd
}

func PreCmdPostToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}

func CmdPostToken(c *core.CommandConfig) error {
	if !viper.IsSet(constants.ArgNoHeaders) {
		viper.Set(constants.ArgNoHeaders, true)
	}

	var err error

	id, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	tokenPostProperties.SetName(name)

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

		tokenPostProperties.SetExpiryDate(expiryDate)
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
		tokenPostProperties.SetExpiryDate(timeNow)
	}

	if viper.IsSet(core.GetFlagName(c.NS, FlagStatus)) {
		var status string

		status, err = c.Command.Command.Flags().GetString(FlagStatus)
		if err != nil {
			return err
		}

		tokenPostProperties.SetStatus(status)
	}

	tokenInput := containerregistry.NewPostTokenInputWithDefaults()
	tokenInput.SetProperties(*tokenPostProperties)

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensPost(context.Background(), id).PostTokenInput(*tokenInput).Execute()
	if err != nil {
		return err
	}

	tokenPrint := containerregistry.NewTokenResponseWithDefaults()
	tokenPrint.SetProperties(token.GetProperties())

	cols := c.Cols()
	// Default to showing only CredentialsPassword for create operations
	if cols == nil {
		cols = []string{"CredentialsPassword"}
	}
	return c.Out(table.Sprint(allCols, token, cols))
}
