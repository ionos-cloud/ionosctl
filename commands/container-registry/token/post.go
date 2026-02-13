package token

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
			ShortDesc: "Create a new token",
			LongDesc:  "Create a new token used to access a container registry",
			Example: "ionosctl container-registry token create --registry-id [REGISTRY-ID] --name [TOKEN-NAME]\n" +
				"In order to save the token to a environment variable: export [ENV-VAL-NAME]=$(ionosctl cr token create --name [TOKEN-NAME] --registry-id [REGISTRY-ID]",
			PreCmdRun:  PreCmdPostToken,
			CmdRun:     CmdPostToken,
			InitClient: true,
		},
	)

	// This line is only used to override the help text for `--no-headers`!
	cmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, true, "Use --no-headers=false to show column headers",
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the Token", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagExpiryDate, "", "", "Expiry date of the Token")
	cmd.AddStringFlag(FlagStatus, "", "", "Status of the Token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"enabled", "disabled",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTimeUntilExpiry, "", "", "Time until the Token expires (ex: 1y2d)")
	cmd.AddStringFlag(
		constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
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

	if c.Command.Command.Flags().Changed(FlagExpiryDate) {
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
	} else if c.Command.Command.Flags().Changed(FlagTimeUntilExpiry) {
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

	if c.Command.Command.Flags().Changed(FlagStatus) {
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

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.ContainerRegistryToken, token, tabheaders.GetHeaders(AllTokenCols, postHeaders, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
