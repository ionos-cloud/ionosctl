package token

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "delete",
			Aliases:    []string{"d", "del", "rm"},
			ShortDesc:  "Delete a token",
			LongDesc:   "Delete a token from a registry",
			Example:    "ionosctl container-registry token delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]",
			PreCmdRun:  PreCmdDeleteToken,
			CmdRun:     CmdDeleteToken,
			InitClient: true,
		},
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all tokens from all registries")
	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddBoolFlag(FlagAllTokens, "", false, "Delete all tokens from a registry")
	cmd.AddStringFlag(FlagTokenId, "", "", "Token ID")
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

	return cmd
}

func CmdDeleteToken(c *core.CommandConfig) error {
	allFlag := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))

	if !allFlag {
		regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
		allTokensFlag := viper.GetBool(core.GetFlagName(c.NS, FlagAllTokens))

		if !allTokensFlag {
			tokenId := viper.GetString(core.GetFlagName(c.NS, FlagTokenId))

			msg := fmt.Sprintf("delete Token: %s", tokenId)

			if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err := client.Must().RegistryClient.TokensApi.RegistriesTokensDelete(context.Background(), regId, tokenId).Execute()
			if err != nil {
				return err
			}

			return nil
		}

		// delete all tokens of a single registry
		return core.DeleteAll(c, core.DeleteAllOptions[containerregistry.TokenResponse]{
			Resource: "Token",
			List: func() ([]containerregistry.TokenResponse, error) {
				tokens, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), regId).Execute()
				if err != nil {
					return nil, err
				}
				return tokens.GetItems(), nil
			},
			Summary: func(t containerregistry.TokenResponse) string {
				return fmt.Sprintf("name: %s, id: %s, registry: %s", t.Properties.Name, *t.Id, regId)
			},
			ID: func(t containerregistry.TokenResponse) string {
				return *t.Id
			},
			Delete: func(t containerregistry.TokenResponse) error {
				_, err := client.Must().RegistryClient.TokensApi.RegistriesTokensDelete(context.Background(), regId, *t.Id).Execute()
				return err
			},
		})
	}

	// delete all tokens across all registries
	return core.DeleteAll(c, core.DeleteAllOptions[regToken]{
		Resource: "Token",
		List: func() ([]regToken, error) {
			regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
			if err != nil {
				return nil, err
			}

			var all []regToken
			for _, reg := range regs.GetItems() {
				tokens, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), *reg.Id).Execute()
				if err != nil {
					return nil, err
				}
				for _, tok := range tokens.GetItems() {
					all = append(all, regToken{regID: *reg.Id, tok: tok})
				}
			}
			return all, nil
		},
		Summary: func(rt regToken) string {
			return fmt.Sprintf("name: %s, id: %s, registry: %s", rt.tok.Properties.Name, *rt.tok.Id, rt.regID)
		},
		ID: func(rt regToken) string {
			return *rt.tok.Id
		},
		Delete: func(rt regToken) error {
			_, err := client.Must().RegistryClient.TokensApi.RegistriesTokensDelete(context.Background(), rt.regID, *rt.tok.Id).Execute()
			return err
		},
	})
}

// regToken pairs a token with its owning registry id so that tokens from all
// registries can be flattened into a single DeleteAll slice.
type regToken struct {
	regID string
	tok   containerregistry.TokenResponse
}

func PreCmdDeleteToken(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagRegistryId, FlagTokenId},
		[]string{constants.FlagRegistryId, FlagAllTokens},
		[]string{constants.ArgAll},
	)
}
