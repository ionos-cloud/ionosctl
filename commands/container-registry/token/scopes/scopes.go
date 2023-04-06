package scopes

import (
	"context"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenScopesCmd() *core.Command {
	scopesCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "scope",
			Aliases: []string{"s", "scopes"},
			Short:   "Registry Tokens' Scopes Operations",
			Long: "Manage token scopes for a registry. You can list, add, and remove scopes from a token. " +
				"Scopes are used to grant access to a registry. " +
				"Each token can have multiple scopes. ",
			TraverseChildren: true,
		},
	}

	scopesCmd.AddCommand(TokenScopesListCmd())
	scopesCmd.AddCommand(TokenScopesAddCmd())
	scopesCmd.AddCommand(TokenScopesDeleteCmd())
	return scopesCmd
}

func getTokenScopePrint(
	resp *shared.APIResponse, c *core.CommandConfig, response *sdkgo.TokenResponse,
	first bool,
) printer.Result {
	r := printer.Result{}

	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(
				core.GetFlagName(
					c.NS, constants.ArgWaitForRequest,
				),
			) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if response != nil {
			r.OutputJSON = response
			r.Columns = printer.GetHeadersListAll(allColsScopes, defaultScopeCols, "TokenId", nil, viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))) // headers
			r.KeyValue = getTokensScopeRows(response)                                                                                                       // map header -> rows
		}
	}
	return r
}

type TokenScopePrint struct {
	ScopeId     string `json:"ScopeId,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
	TokenId     string `json:"TokenId,omitempty"`
	Type        string `json:"Type,omitempty"`
	Actions     string `json:"Actions,omitempty"`
}

func getTokensScopeRows(token *sdkgo.TokenResponse) []map[string]interface{} {
	scopes := token.Properties.GetScopes()
	out := make([]map[string]interface{}, 0, len(*scopes))
	for i, scope := range *scopes {
		var tokenScopePrint TokenScopePrint
		if nameOk, ok := scope.GetNameOk(); ok && nameOk != nil {
			tokenScopePrint.DisplayName = *nameOk
		}
		if typeOk, ok := scope.GetTypeOk(); ok && typeOk != nil {
			tokenScopePrint.Type = *typeOk
		}
		if actionsOk, ok := scope.GetActionsOk(); ok && actionsOk != nil {
			for _, action := range *actionsOk {
				tokenScopePrint.Actions += string(action) + ", "
			}
			tokenScopePrint.Actions = tokenScopePrint.Actions[:len(tokenScopePrint.Actions)-2]
		}
		tokenScopePrint.TokenId = *token.Id
		tokenScopePrint.ScopeId = strconv.Itoa(i + 1)
		o := structs.Map(tokenScopePrint)
		out = append(out, o)
	}
	return out
}

var allColsScopes = structs.Names(TokenScopePrint{})

func TokensIds(regId string) []string {
	svcToken := resources.NewTokenService(client.Must(), context.Background())
	var allTokens []sdkgo.TokenResponse

	if regId != "" {

		tokens, _, _ := svcToken.List(regId)

		allTokens = append(allTokens, *tokens.GetItems()...)

		return shared.Map(
			allTokens, func(reg sdkgo.TokenResponse) string {
				return *reg.GetId()
			},
		)
	}

	svc := resources.NewRegistriesService(client.Must(), context.Background())
	regs, _, _ := svc.List("")
	regsIDs := *regs.GetItems()

	for _, regID := range regsIDs {
		tokens, _, _ := svcToken.List(*regID.GetId())

		allTokens = append(allTokens, *tokens.GetItems()...)
	}

	return shared.Map(
		allTokens, func(reg sdkgo.TokenResponse) string {
			return *reg.GetId()
		},
	)
}

var (
	defaultScopeCols = []string{"ScopeId", "DisplayName", "Type", "Actions"}
	allScopeCols     = []string{"ScopeId", "TokenId", "DisplayName", "Type", "Actions"}
)
