package scopes

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var allScopeCols = []table.Column{
	{Name: "ScopeId", Default: true, Format: func(item map[string]any) any {
		// ScopeId is the index; set externally via SetCell
		return item["ScopeId"]
	}},
	{Name: "DisplayName", JSONPath: "name", Default: true},
	{Name: "Type", JSONPath: "type", Default: true},
	{Name: "Actions", Default: true, Format: func(item map[string]any) any {
		actions := table.Navigate(item, "actions")
		if actions == nil {
			return nil
		}
		arr, ok := actions.([]any)
		if !ok {
			return fmt.Sprintf("%v", actions)
		}
		parts := make([]string, len(arr))
		for i, a := range arr {
			parts[i] = fmt.Sprintf("%v", a)
		}
		return strings.Join(parts, ", ")
	}},
}

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

	scopesCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allScopeCols))
	_ = scopesCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allScopeCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	scopesCmd.AddCommand(TokenScopesListCmd())
	scopesCmd.AddCommand(TokenScopesAddCmd())
	scopesCmd.AddCommand(TokenScopesDeleteCmd())
	return scopesCmd
}

func TokensIds(regId string) []string {
	var allTokens []containerregistry.TokenResponse

	if regId != "" {
		// list tokens for provided registry
		tokens, _, _ := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), regId).Execute()
		allTokens = append(allTokens, tokens.GetItems()...)
		return functional.Map(allTokens, func(t containerregistry.TokenResponse) string { return t.GetId() })
	}

	// list all registries then tokens for each
	regs, _, _ := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
	for _, reg := range regs.GetItems() {
		toks, _, _ := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), reg.GetId()).Execute()
		allTokens = append(allTokens, toks.GetItems()...)
	}
	return functional.Map(allTokens, func(t containerregistry.TokenResponse) string { return t.GetId() })
}
