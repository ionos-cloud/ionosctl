package authv1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/auth-v1/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	authv1 "github.com/ionos-cloud/ionosctl/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/services/auth-v1/resources"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenCmd() *core.Command {
	ctx := context.TODO()
	tokenCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "token",
			Short:            "Token Operations",
			Long:             "The sub-commands of `ionosctl token` allow you to generate, list, get, delete Tokens.",
			TraverseChildren: true,
		},
	}
	globalFlags := tokenCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultTokenCols, printer.ColsMessage(allTokenCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(tokenCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = tokenCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allTokenCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Tokens",
		LongDesc:   "Use this command to retrieve a complete list of Tokens under your account, to list the details of all tokens, generated by the user.",
		Example:    listTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTokenList,
		InitClient: true,
	})
	list.AddIntFlag(authv1.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are listed")
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specified Token",
		LongDesc:   "Use this command to retrieve details about a Token by using its ID.\n\nRequired values to run command:\n\n* Token Id",
		Example:    getTokenExample,
		PreCmdRun:  PreRunTokenId,
		CmdRun:     RunTokenGet,
		InitClient: true,
	})
	get.AddStringFlag(authv1.ArgTokenId, authv1.ArgIdShort, "", authv1.TokenId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(authv1.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddIntFlag(authv1.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the token information is displayed")
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Generate/Create Command
	*/
	generate := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "generate",
		Aliases:    []string{"create"},
		ShortDesc:  "Create a new Token",
		LongDesc:   "Use this command to generate a new Token. Only the JSON Web Token, associated with user credentials, will be displayed.",
		Example:    generateTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTokenCreate,
		InitClient: true,
	})
	generate.AddIntFlag(authv1.ArgContractNo, "", 0, "Users with multiple contracts can provide the contract number, for which the token is generated")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace: "token",
		Resource:  "token",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete one or multiple Tokens",
		LongDesc: `Use this command to delete a specified Token by token Id or multiple Tokens (based on a criteria: CURRENT, EXPIRED, ALL) from your account. With parameter values ALL and EXPIRED, 'Basic Authentication' or 'Token Authentication' tokens with valid credentials must be encapsulated in the header. With value CURRENT, only the 'Token Authentication' with valid credentials is required.

Required values to run command:

* Token Id/CURRENT/EXPIRED/ALL`,
		Example:    deleteTokenExample,
		PreCmdRun:  PreRunTokenDelete,
		CmdRun:     RunTokenDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(authv1.ArgTokenId, authv1.ArgIdShort, "", authv1.TokenId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(authv1.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(authv1.ArgCurrent, authv1.ArgCurrentShort, false, "Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file", core.RequiredFlagOption())
	deleteCmd.AddBoolFlag(authv1.ArgExpired, authv1.ArgExpiredShort, false, "Delete the Tokens that are currently expired", core.RequiredFlagOption())
	deleteCmd.AddBoolFlag(authv1.ArgAll, authv1.ArgAllShort, false, "Delete the Tokens under your account", core.RequiredFlagOption())
	deleteCmd.AddIntFlag(authv1.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are deleted")

	return tokenCmd
}

const contractNumberMessage = "Contract Number: %v"

func PreRunTokenId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, authv1.ArgTokenId)
}

func PreRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authv1.ArgTokenId}, []string{authv1.ArgCurrent}, []string{authv1.ArgExpired}, []string{authv1.ArgAll})
}

func RunTokenList(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	}
	tokens, _, err := c.AuthV1Services.Tokens().List(viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, getTokens(tokens)))
}

func RunTokenGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Token with ID: %v...", viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)))
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	}
	token, _, err := c.AuthV1Services.Tokens().Get(viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)), viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, []resources.Token{*token}))
}

func RunTokenCreate(c *core.CommandConfig) error {
	c.Printer.Verbose("Generating new token..")
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	}
	newJwt, _, err := c.AuthV1Services.Tokens().Create(viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if newJwt != nil {
		if token, ok := newJwt.GetTokenOk(); ok && token != nil {
			return c.Printer.Print(*token)
		} else {
			return errors.New("error getting generated token")
		}
	} else {
		return errors.New("error getting generated JWT")
	}
}

func RunTokenDelete(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgTokenId)) {
		return RunTokenDeleteById(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgCurrent)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgCurrent)) {
		return RunTokenDeleteCurrent(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgExpired)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgExpired)) {
		return RunTokenDeleteExpired(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgAll)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgAll)) {
		return RunTokenDeleteAll(c)
	}
	return nil
}

func RunTokenDeleteAll(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete all tokens")); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting all tokens...")
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("ALL", viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Print("Status: all tokens have been successfully deleted")
			}
		}
	}
	return errors.New("error deleting all tokens")
}

func RunTokenDeleteExpired(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete expired tokens")); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting expired tokens...")
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("EXPIRED", viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Print("Status: expired tokens have been successfully deleted")
			}
		}
	}
	return errors.New("error deleting expired tokens")
}

func RunTokenDeleteCurrent(c *core.CommandConfig) error {
	c.Printer.Verbose("Note: This operation is based on Authorization Header for Bearer Token")
	if viper.GetString(config.Token) == "" {
		return errors.New(fmt.Sprintf("no token found. Please make sure you have exported the %s environment variable or you have token set in the config file",
			sdkgoauth.IonosTokenEnvVar))
	}
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete CURRENT token")); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting current token...")
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("CURRENT", viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Print("Status: current token have been successfully deleted")
			}
		}
	}
	return errors.New("error deleting current token")
}

func RunTokenDeleteById(c *core.CommandConfig) error {
	tokenId := viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId))
	c.Printer.Verbose("Token ID: %s", tokenId)
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete token with ID: %s", tokenId)); err != nil {
		return err
	}
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(tokenId, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Print("Status: token has been successfully deleted")
			}
		}
	}
	return errors.New("error deleting token")
}

// Output Printing

var (
	defaultTokenCols = []string{"TokenId", "CreatedDate", "ExpirationDate"}
	allTokenCols     = []string{"TokenId", "CreatedDate", "ExpirationDate", "Href"}
)

type TokenPrint struct {
	TokenId        string `json:"TokenId,omitempty"`
	CreatedDate    string `json:"CreatedDate,omitempty"`
	ExpirationDate string `json:"ExpirationDate,omitempty"`
	Href           string `json:"Href,omitempty"`
}

func getTokenPrint(c *core.CommandConfig, dcs []resources.Token) printer.Result {
	r := printer.Result{}
	if c != nil {
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getTokensKVMaps(dcs)
			r.Columns = getTokenCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getTokenCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultTokenCols
	}

	columnsMap := map[string]string{
		"TokenId":        "TokenId",
		"CreatedDate":    "CreatedDate",
		"ExpirationDate": "ExpirationDate",
		"Href":           "Href",
	}
	var tokenCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			tokenCols = append(tokenCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return tokenCols
}

func getTokens(tokens resources.Tokens) []resources.Token {
	dc := make([]resources.Token, 0)
	if items, ok := tokens.GetTokensOk(); ok && items != nil {
		for _, token := range *items {
			dc = append(dc, resources.Token{Token: token})
		}
	}
	return dc
}

func getTokensKVMaps(tokens []resources.Token) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(tokens))
	for _, token := range tokens {
		var tokenPrint TokenPrint
		if idOk, ok := token.GetIdOk(); ok && idOk != nil {
			tokenPrint.TokenId = *idOk
		}
		if createdDateOk, ok := token.GetCreatedDateOk(); ok && createdDateOk != nil {
			tokenPrint.CreatedDate = *createdDateOk
		}
		if expirationDateOk, ok := token.GetExpirationDateOk(); ok && expirationDateOk != nil {
			tokenPrint.ExpirationDate = *expirationDateOk
		}
		if hrefOk, ok := token.GetHrefOk(); ok && hrefOk != nil {
			tokenPrint.Href = *hrefOk
		}
		o := structs.Map(tokenPrint)
		out = append(out, o)
	}
	return out
}
