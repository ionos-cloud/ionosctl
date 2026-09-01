package scopes

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenScopesAddCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "token",
			Resource:  "scope",
			Verb:      "add",
			Aliases:   []string{"a", "ad"},
			ShortDesc: "Add an access scope to a token",
			LongDesc: `Add a scope to an existing token, granting it a set of actions on a target resource. The scope is appended to the token's existing scopes (the token is patched, not replaced), so its password and other scopes are preserved.

--name is the target resource (a repository name/path, or '*' for all), --type is the kind of target ('repository', 'namespace' or 'registry'), and --actions is the comma-separated list of allowed operations ('pull', 'push', 'delete', or '*').`,
			Example: `# Grant pull+push on a single repository
ionosctl container-registry token scope add --registry-id REGISTRY_ID --token-id TOKEN_ID --name my-app --type repository --actions pull,push

# Grant read-only (pull) access to every repository in the registry
ionosctl container-registry token scope add --registry-id REGISTRY_ID --token-id TOKEN_ID --name "*" --type repository --actions pull`,
			PreCmdRun:  PreCmdTokenScopesAdd,
			CmdRun:     CmdTokenScopesAdd,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry that owns the token", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "", "", "The unique ID of the token to add the scope to")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId,
		func(cobracmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(
				viper.GetString(
					core.GetFlagName(
						cmd.NS, constants.FlagRegistryId,
					),
				),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Target resource of the scope: a repository name/path, or '*' for all repositories", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagType, "y", "", "Kind of target the --name refers to: 'repository' (one repo), 'namespace' (a repo path prefix) or 'registry' (the whole registry)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagType,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"repository", "namespace", "registry"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringSliceFlag(FlagActions, "a", []string{}, "Comma-separated operations the token may perform on the target, e.g. pull, push, delete (or '*' for all)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagActions,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"*", "push", "pull", "delete", "read", "write", "list"}, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func PreCmdTokenScopesAdd(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagRegistryId, FlagTokenId, constants.FlagName, FlagActions, FlagType,
	)
	if err != nil {
		return err
	}
	return nil
}

func CmdTokenScopesAdd(c *core.CommandConfig) error {
	var scope containerregistry.Scope
	var err error

	regId, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	tokenId, err := c.Command.Command.Flags().GetString(FlagTokenId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	actions, err := c.Command.Command.Flags().GetStringSlice(FlagActions)
	if err != nil {
		return err
	}

	scopeType, err := c.Command.Command.Flags().GetString(FlagType)
	if err != nil {
		return err
	}

	scope.SetName(name)
	scope.SetActions(actions)
	scope.SetType(scopeType)

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensFindById(context.Background(), regId, tokenId).Execute()
	if err != nil {
		return err
	}

	updateToken := containerregistry.NewPatchTokenInput()
	if token.Properties.ExpiryDate != nil {
		updateToken.SetExpiryDate(token.Properties.GetExpiryDate())
	}

	if token.Properties.Status != nil {
		updateToken.SetStatus(token.Properties.GetStatus())
	}

	scopes := token.Properties.GetScopes()
	scopes = append(scopes, scope)

	updateToken.SetScopes(scopes)

	tokenUp, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensPatch(context.Background(), regId, tokenId).PatchTokenInput(*updateToken).Execute()
	if err != nil {
		return err
	}

	t := table.New(allScopeCols)
	if err := t.Extract(tokenUp.Properties.Scopes); err != nil {
		return err
	}
	for i := range tokenUp.Properties.Scopes {
		t.SetCell(i, "ScopeId", fmt.Sprintf("%v", i))
	}

	return c.Out(t.Render(table.ResolveCols(allScopeCols, c.Cols())))
}
