package authv1

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenCmd() *core.Command {
	tokenCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "token",
			Short:            "Token Operations",
			Long:             "The sub-commands of `ionosctl token` allow you to generate, list, get, delete Tokens.",
			TraverseChildren: true,
		},
	}
	globalFlags := tokenCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", nil, printer.ColsMessage(allTokenCols))
	_ = viper.BindPFlag(core.GetFlagName(tokenCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = tokenCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allTokenCols, cobra.ShellCompDirectiveNoFileComp
	})

	tokenCmd.AddCommand(TokenPostCmd())
	tokenCmd.AddCommand(TokenListCmd())
	tokenCmd.AddCommand(TokenGetCmd())
	tokenCmd.AddCommand(TokenDeleteCmd())
	tokenCmd.AddCommand(TokenParseCmd())

	return tokenCmd
}

const contractNumberMessage = "Contract Number: %v"

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
			r.Columns = printer.GetHeadersAllDefault(allTokenCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)))
		}
	}
	return r
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
