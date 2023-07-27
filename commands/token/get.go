package token

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/token/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
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
		return completer.TokensIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(authservice.ArgToken, authservice.ArgTokenShort, "", authservice.Token, core.RequiredFlagOption())
	cmd.AddIntFlag(authservice.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the token information is displayed")
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
	c.Printer.Verbose("Getting Token with ID: %v...", viper.GetString(core.GetFlagName(c.NS, authservice.ArgTokenId)))
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	}
	token, _, err := c.AuthV1Services.Tokens().Get(viper.GetString(core.GetFlagName(c.NS, authservice.ArgTokenId)), viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, []resources.Token{*token}))
}

func runTokenGetByToken(c *core.CommandConfig) error {
	token := viper.GetString(core.GetFlagName(c.NS, authservice.ArgToken))
	c.Printer.Verbose("Token content is: %s", token)

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return err
	}

	c.Printer.Verbose("Getting Token with ID: %v...", tokenId)
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	}
	tokenObj, _, err := c.AuthV1Services.Tokens().Get(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, []resources.Token{*tokenObj}))
}
