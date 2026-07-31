package scopes

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
			Short:   "Manage a token's access scopes",
			Long: `A scope grants a token a set of actions on a target resource. A token with no scopes cannot pull or push anything; adding scopes is what makes a token usable.

Each scope has three parts:
  - name: the resource the scope applies to, e.g. a repository name/path, or '*' for all repositories;
  - type: the kind of resource the name refers to - typically 'repository' (a single repo), 'namespace' (a repo path prefix), or 'registry' (the whole registry);
  - actions: the operations allowed, e.g. 'pull', 'push', 'delete' (or '*' for all).

A token may hold multiple scopes. In this CLI each scope is addressed by its zero-based index within the token (ScopeId), shown by 'scope list'; adding a scope appends to the list, deleting removes by that index.`,
			TraverseChildren: true,
		},
	}

	scopesCmd.AddColsFlag(allScopeCols)

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
