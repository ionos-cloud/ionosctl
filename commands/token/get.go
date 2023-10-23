package token

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/token/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specified Token",
		LongDesc:   "Use this command to retrieve details about a Token by using its ID.\n\nRequired values to run command:\n\n* Token Id",
		Example:    getTokenExample,
		PreCmdRun:  preRunTokenId,
		CmdRun:     runTokenGet,
		InitClient: true,
	})

	cmd.AddUUIDFlag(authservice.ArgTokenId, authservice.ArgIdShort, "", authservice.TokenId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(authservice.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(authservice.ArgToken, authservice.ArgTokenShort, "", authservice.Token, core.RequiredFlagOption())
	cmd.AddIntFlag(authservice.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the token information is displayed")

	return cmd
}

func preRunTokenId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authservice.ArgTokenId}, []string{authservice.ArgToken})
}

func runTokenGet(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgTokenId)) {
		return runTokenGetById(c)
	}

	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgToken)) {
		return runTokenGetByToken(c)
	}

	return nil
}

func runTokenGetById(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting Token with ID: %v...", viper.GetString(core.GetFlagName(c.NS, authservice.ArgTokenId))))

	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))))
	}

	token, _, err := c.AuthV1Services.Tokens().Get(viper.GetString(core.GetFlagName(c.NS, authservice.ArgTokenId)), viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.AuthToken, token.Token,
		tabheaders.GetHeaders(allTokenCols, defaultTokenCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func runTokenGetByToken(c *core.CommandConfig) error {
	token := viper.GetString(core.GetFlagName(c.NS, authservice.ArgToken))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Token content is: %s", token))

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Token with ID: %v...", tokenId))
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))))
	}
	tokenObj, _, err := c.AuthV1Services.Tokens().Get(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.AuthToken, tokenObj.Token,
		tabheaders.GetHeaders(allTokenCols, defaultTokenCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
