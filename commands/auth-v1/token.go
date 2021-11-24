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
	core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Tokens",
		LongDesc:   "Use this command to retrieve a complete list of Tokens under your account.",
		Example:    listTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTokenList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Token ",
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

	/*
		Generate/Create Command
	*/
	core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "generate",
		Aliases:    []string{"g"},
		ShortDesc:  "Create a new Token",
		LongDesc:   "Use this command to generate a new Token. Only the JWT will be displayed.",
		Example:    generateTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTokenCreate,
		InitClient: true,
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace: "token",
		Resource:  "token",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Token",
		LongDesc: `Use this command to delete a specified Token or multiple Tokens (based on a criteria: CURRENT, EXPIRED, ALL) from your account.

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

	return tokenCmd
}

func PreRunTokenId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, authv1.ArgTokenId)
}

func PreRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authv1.ArgTokenId}, []string{authv1.ArgCurrent}, []string{authv1.ArgExpired}, []string{authv1.ArgAll})
}

func RunTokenList(c *core.CommandConfig) error {
	tokens, _, err := c.AuthV1Services.Tokens().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, getTokens(tokens)))
}

func RunTokenGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Token with ID: %v...", viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)))
	token, _, err := c.AuthV1Services.Tokens().Get(viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, []resources.Token{*token}))
}

func RunTokenCreate(c *core.CommandConfig) error {
	newJwt, _, err := c.AuthV1Services.Tokens().Create()
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
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgTokenId)) {
		tokenId := viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId))
		c.Printer.Verbose("Token ID: %s", tokenId)
		if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete token with ID: %s", tokenId)); err != nil {
			return err
		}
		tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(tokenId)
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
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgCurrent)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgCurrent)) {
		c.Printer.Verbose("Note: Authentication based on token needs to be used in order for the deletion to succeed.")
		if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete CURRENT token")); err != nil {
			return err
		}
		tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("CURRENT")
		if err != nil {
			return err
		}
		if tokenResponse != nil {
			if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
				if *success {
					return c.Printer.Print("Status: tokens based on criteria: CURRENT have been successfully deleted")
				}
			}
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgExpired)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgExpired)) {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete expired tokens")); err != nil {
			return err
		}
		tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("EXPIRED")
		if err != nil {
			return err
		}
		if tokenResponse != nil {
			if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
				if *success {
					return c.Printer.Print("Status: tokens based on criteria: EXPIRED have been successfully deleted")
				}
			}
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgAll)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgAll)) {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete all tokens")); err != nil {
			return err
		}
		tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("ALL")
		if err != nil {
			return err
		}
		if tokenResponse != nil {
			if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
				if *success {
					return c.Printer.Print("Status: tokens based on criteria: ALL have been successfully deleted")
				}
			}
		}
	}
	return c.Printer.Print("Status: token delete command has been successfully executed")
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
