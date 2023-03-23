//go:build integration
// +build integration

package token

import (
	"context"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/lucasjones/reggen"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTokenService(t *testing.T) {
	t.Run("token functions", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)

		// create registry
		name, err := reggen.Generate("^[a-z][-a-z0-9]{1,61}[a-z0-9]$", 10)
		assert.NoError(t, err)
		c := registry.RegPostCmd()
		c.Command.Flags().Set(FlagName, name)
		c.Command.Flags().Set(registry.FlagLocation, "de/fra")

		err = c.Command.Execute()
		assert.NoError(t, err)
		time.Sleep(10 * time.Second)

		// get registry
		assert.NoError(t, err)
		registries, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
		assert.NoError(t, err)

		var newReg *ionoscloud.RegistryResponse
		for _, reg := range *registries.GetItems() {
			if *reg.GetProperties().GetName() == name {
				newReg = &reg
			}
		}

		time.Sleep(10 * time.Second)
		// create token
		tokenName, err := reggen.Generate("^[A-Za-z][-A-Za-z0-9]{0,61}[A-Za-z0-9]$", 10)
		cToken := TokenPostCmd()
		cToken.Command.Flags().Set(FlagRegId, *newReg.GetId())
		cToken.Command.Flags().Set(FlagName, tokenName)
		assert.NoError(t, err)

		err = cToken.Command.Execute()
		assert.NoError(t, err)

		time.Sleep(10 * time.Second)
		// get token
		tokens, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), *newReg.GetId()).Execute()
		assert.NoError(t, err)

		var newToken *ionoscloud.TokenResponse
		for _, token := range *tokens.GetItems() {
			if *token.GetProperties().GetName() == tokenName {
				newToken = &token
			}
		}

		time.Sleep(10 * time.Second)
		// add scopes
		addScopes := scopes.TokenScopesAddCmd()
		addScopes.Command.Flags().Set(FlagRegId, *newReg.GetId())
		addScopes.Command.Flags().Set(FlagTokenId, *newToken.GetId())
		addScopes.Command.Flags().Set(scopes.FlagName, "test")
		addScopes.Command.Flags().Set(scopes.FlagType, "repository")
		addScopes.Command.Flags().Set(scopes.FlagActions, "pull")

		err = addScopes.Command.Execute()
		assert.NoError(t, err)
		time.Sleep(10 * time.Second)

		addScopes = scopes.TokenScopesAddCmd()
		addScopes.Command.Flags().Set(FlagRegId, *newReg.GetId())
		addScopes.Command.Flags().Set(FlagTokenId, *newToken.GetId())
		addScopes.Command.Flags().Set(scopes.FlagName, "test2")
		addScopes.Command.Flags().Set(scopes.FlagType, "repository")
		addScopes.Command.Flags().Set(scopes.FlagActions, "push")
		time.Sleep(10 * time.Second)

		err = addScopes.Command.Execute()
		assert.NoError(t, err)

		// delete scopes
		deleteScopes := scopes.TokenScopesDeleteCmd()
		deleteScopes.Command.Flags().Set(FlagRegId, *newReg.GetId())
		deleteScopes.Command.Flags().Set(FlagTokenId, *newToken.GetId())
		deleteScopes.Command.Flags().Set(constants.ArgAll, "true")

		err = deleteScopes.Command.Execute()
		assert.NoError(t, err)

		time.Sleep(10 * time.Second)

		gToken := TokenGetCmd()
		gToken.Command.Flags().Set(FlagRegId, *newReg.GetId())
		gToken.Command.Flags().Set(FlagTokenId, *newToken.GetId())

		err = gToken.Command.Execute()
		assert.NoError(t, err)

		uToken := TokenUpdateCmd()
		uToken.Command.Flags().Set(FlagRegId, *newReg.GetId())
		uToken.Command.Flags().Set(FlagTokenId, *newToken.GetId())
		uToken.Command.Flags().Set(FlagName, "newName")
		uToken.Command.Flags().Set(FlagStatus, "disabled")

		err = uToken.Command.Execute()
		assert.NoError(t, err)
		time.Sleep(10 * time.Second)

		checkProp, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensFindById(context.Background(), *newReg.GetId(), *newToken.GetId()).Execute()
		assert.NoError(t, err)
		assert.Equal(t, "disabled", *checkProp.GetProperties().GetStatus())
		time.Sleep(10 * time.Second)
		// delete token
		dToken := TokenDeleteCmd()
		dToken.Command.Flags().Set(FlagRegId, *newReg.GetId())
		dToken.Command.Flags().Set(FlagTokenId, *newToken.GetId())

		err = dToken.Command.Execute()
		assert.NoError(t, err)
		time.Sleep(10 * time.Second)

		// replace token
		rToken := TokenReplaceCmd()
		rToken.Command.Flags().Set(FlagRegId, *newReg.GetId())
		rToken.Command.Flags().Set(FlagTokenId, *newToken.GetId())
		rToken.Command.Flags().Set(FlagName, "newName")
		time.Sleep(10 * time.Second)

		err = rToken.Command.Execute()
		assert.NoError(t, err)
		time.Sleep(10 * time.Second)

		// delete registry
		d := registry.RegDeleteCmd()
		d.Command.Flags().Set(FlagRegId, *newReg.GetId())

		err = d.Command.Execute()
		assert.NoError(t, err)

	})
}
