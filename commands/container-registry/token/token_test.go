package token

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/lucasjones/reggen"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistryService(t *testing.T) {
	t.Run("registry functions", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)

		name, err := reggen.Generate("^[a-z][-a-z0-9]{1,61}[a-z0-9]$", 10)
		assert.NoError(t, err)
		c := registry.RegPostCmd()
		c.Command.Flags().Set(FlagName, name)
		c.Command.Flags().Set(registry.FlagLocation, "de/fra")

		err = c.Command.Execute()
		assert.NoError(t, err)

		svc, err := config.GetClient()
		registries, _, err := svc.RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
		assert.NoError(t, err)

		var newReg *ionoscloud.RegistryResponse
		for _, reg := range *registries.GetItems() {
			if *reg.GetProperties().GetName() == name {
				newReg = &reg
			}
		}

		tokenName, err := reggen.Generate("^[A-Za-z][-A-Za-z0-9]{0,61}[A-Za-z0-9]$", 10)
		cToken := TokenPostCmd()
		cToken.Command.Flags().Set(FlagRegId, *newReg.GetId())
		cToken.Command.Flags().Set(FlagName, tokenName)
		assert.NoError(t, err)

		err = cToken.Command.Execute()
		assert.NoError(t, err)

		tokens, _, err := svc.RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), *newReg.GetId()).Execute()
		assert.NoError(t, err)

		var newToken *ionoscloud.TokenResponse
		for _, token := range *tokens.GetItems() {
			if *token.GetProperties().GetName() == tokenName {
				newToken = &token
			}
		}

		addScopes := scopes.TokenScopesAddCmd()
		addScopes.Command.Flags().Set(FlagRegId, *newReg.GetId())
		addScopes.Command.Flags().Set(FlagTokenId, *newToken.GetId())
		addScopes.Command.Flags().Set(scopes.FlagName, "test")
		addScopes.Command.Flags().Set(scopes.FlagType, "repository")
		addScopes.Command.Flags().Set(scopes.FlagActions, "pull")

		err = addScopes.Command.Execute()
		assert.NoError(t, err)

	})
}
