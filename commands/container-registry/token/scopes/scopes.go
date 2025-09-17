package scopes

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var (
	allScopeCols = []string{"ScopeId", "DisplayName", "Type", "Actions"}
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
