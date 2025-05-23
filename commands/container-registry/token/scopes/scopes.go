package scopes

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"

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
	svcToken := resources.NewTokenService(client.Must(), context.Background())
	var allTokens []containerregistry.TokenResponse

	if regId != "" {
		tokens, _, _ := svcToken.List(regId)

		allTokens = append(allTokens, tokens.GetItems()...)

		return functional.Map(
			allTokens, func(reg containerregistry.TokenResponse) string {
				return reg.GetId()
			},
		)
	}

	svc := resources.NewRegistriesService(client.Must(), context.Background())
	regs, _, _ := svc.List("")
	regsIDs := regs.GetItems()

	for _, regID := range regsIDs {
		tokens, _, _ := svcToken.List(regID.GetId())

		allTokens = append(allTokens, tokens.GetItems()...)
	}

	return functional.Map(
		allTokens, func(reg containerregistry.TokenResponse) string {
			return reg.GetId()
		},
	)
}
