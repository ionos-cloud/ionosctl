package authv1

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"io"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/auth-v1/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
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
		LongDesc:   `Use this command to retrieve a complete list of Tokens under your account.`,
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
		LongDesc:   "Use this command to retrieve details about a Virtual Token  by using its ID.\n\nRequired values to run command:\n\n* Token  Id",
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
	generate := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "generate",
		Aliases:    []string{"g"},
		ShortDesc:  "Create a new Token",
		LongDesc:   `Use this command to generate a new Token.`,
		Example:    generateTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTokenCreate,
		InitClient: true,
	})
	generate.AddBoolFlag(authv1.ArgDecode, authv1.ArgDecodeShort, false, "")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace: "token",
		Resource:  "token",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Token",
		LongDesc: `Use this command to delete a specified Token or multiple Tokens (based on a criteria) from your account.

Required values to run command:

* Token Id/Criteria`,
		Example:    deleteTokenExample,
		PreCmdRun:  PreRunTokenDelete,
		CmdRun:     RunTokenDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(authv1.ArgTokenId, authv1.ArgIdShort, "", authv1.TokenId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(authv1.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(authv1.ArgCriteria, authv1.ArgCriteriaShort, "", "Delete all the Tokens based on this criteria", core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(authv1.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"EXPIRED", "CURRENT", "ALL"}, cobra.ShellCompDirectiveNoFileComp
	})

	return tokenCmd
}

func PreRunTokenId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, authv1.ArgTokenId)
}

func PreRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authv1.ArgTokenId}, []string{authv1.ArgCriteria})
}

func RunTokenList(c *core.CommandConfig) error {
	tokens, _, err := c.AuthV1Services.Tokens().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(nil, c, getTokens(tokens)))
}

func RunTokenGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Token with ID: %v...", viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)))
	dc, _, err := c.AuthV1Services.Tokens().Get(viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(nil, c, []resources.Token{*dc}))
}

func RunTokenCreate(c *core.CommandConfig) error {
	dc, resp, err := c.AuthV1Services.Tokens().Create()
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}

	c.Printer.Print(*dc.Jwt.Token)
	if viper.GetBool(core.GetFlagName(c.NS, authv1.ArgDecode)) {
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "parse jwt token unverified"); err != nil {
			return err
		}
		tokenString := *dc.Jwt.Token
		claims := jwt.MapClaims{}
		parser := &jwt.Parser{}
		token, _, err := parser.ParseUnverified(tokenString, claims)
		if err != nil {
			return err
		}
		c.Printer.Print("Token2:\n")
		c.Printer.Print(token.Header)
		// do something with decoded claims
		for key, val := range claims {
			fmt.Printf("Key: %v, value: %v\n", key, val)
		}
	}
	return nil
}

func RunTokenDelete(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgTokenId)) {
		tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId)))
		if err != nil {
			return err
		}
		return c.Printer.Print(tokenResponse)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgCriteria)) {
		tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria(viper.GetString(core.GetFlagName(c.NS, authv1.ArgCriteria)))
		if err != nil {
			return err
		}
		return c.Printer.Print(tokenResponse)
	}
	return nil
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

func getTokenPrint(resp *resources.Response, c *core.CommandConfig, dcs []resources.Token) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
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
