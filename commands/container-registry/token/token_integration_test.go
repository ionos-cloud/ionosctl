//go:build integration
// +build integration

package token

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func teardown() {
	regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()

	if err != nil {
		log.Print(fmt.Errorf("failed deleting all registries: %w", err))
	}
	for _, reg := range regs.Items {
		_, err := client.Must().RegistryClient.RegistriesApi.RegistriesDelete(context.Background(), *reg.Id).Execute()
		if err != nil {
			log.Print(fmt.Errorf("failed deleting registry: %w", err))
		}
	}
	time.Sleep(30 * time.Second)
}

func TestTokenService(t *testing.T) {
	t.Run(
		"token functions", func(t *testing.T) {
			t.Cleanup(teardown)
			viper.Reset()
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgForce, true)

			// create registry
			name := "ionosctl-crreg-test-" + fake.AlphaNum(8)
			c := registry.RegPostCmd()
			c.Command.Flags().Set(constants.FlagName, name)
			c.Command.Flags().Set(registry.constants.FlagLocation, "de/fra")

			err := c.Command.Execute()
			assert.NoError(t, err)
			time.Sleep(20 * time.Second)

			// get registry
			assert.NoError(t, err)
			registries, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
			assert.NoError(t, err)

			var newReg *containerregistry.RegistryResponse
			for _, reg := range registries.GetItems() {
				if *reg.GetProperties().GetName() == name {
					newReg = &reg
				}
			}

			time.Sleep(20 * time.Second)
			// create token
			tokenName := "ionosctl-crreg-test-" + fake.AlphaNum(8)
			cToken := TokenPostCmd()
			cToken.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			cToken.Command.Flags().Set(constants.FlagName, tokenName)
			assert.NoError(t, err)

			err = cToken.Command.Execute()
			assert.NoError(t, err)

			time.Sleep(20 * time.Second)
			// get token
			tokens, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(
				context.Background(), newReg.GetId(),
			).Execute()
			assert.NoError(t, err)

			var newToken *containerregistry.TokenResponse
			for _, token := range tokens.GetItems() {
				if *token.GetProperties().GetName() == tokenName {
					newToken = &token
				}
			}

			time.Sleep(20 * time.Second)
			// add scopes
			addScopes := scopes.TokenScopesAddCmd()
			addScopes.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			addScopes.Command.Flags().Set(FlagTokenId, newToken.GetId())
			addScopes.Command.Flags().Set(scopes.constants.FlagName, "test")
			addScopes.Command.Flags().Set(scopes.FlagType, "repository")
			addScopes.Command.Flags().Set(scopes.FlagActions, "pull")

			err = addScopes.Command.Execute()
			assert.NoError(t, err)
			time.Sleep(20 * time.Second)

			addScopes = scopes.TokenScopesAddCmd()
			addScopes.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			addScopes.Command.Flags().Set(FlagTokenId, newToken.GetId())
			addScopes.Command.Flags().Set(scopes.constants.FlagName, "test2")
			addScopes.Command.Flags().Set(scopes.FlagType, "repository")
			addScopes.Command.Flags().Set(scopes.FlagActions, "push")
			time.Sleep(20 * time.Second)

			err = addScopes.Command.Execute()
			assert.NoError(t, err)

			// delete scopes
			deleteScopes := scopes.TokenScopesDeleteCmd()
			deleteScopes.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			deleteScopes.Command.Flags().Set(FlagTokenId, newToken.GetId())
			deleteScopes.Command.Flags().Set(constants.ArgAll, "true")

			err = deleteScopes.Command.Execute()
			assert.NoError(t, err)

			time.Sleep(20 * time.Second)

			gToken := TokenGetCmd()
			gToken.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			gToken.Command.Flags().Set(FlagTokenId, newToken.GetId())

			err = gToken.Command.Execute()
			assert.NoError(t, err)

			uToken := TokenUpdateCmd()
			uToken.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			uToken.Command.Flags().Set(FlagTokenId, newToken.GetId())
			uToken.Command.Flags().Set(constants.FlagName, "newName")
			uToken.Command.Flags().Set(FlagStatus, "disabled")

			err = uToken.Command.Execute()
			assert.NoError(t, err)
			time.Sleep(20 * time.Second)

			checkProp, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensFindById(
				context.Background(), newReg.GetId(), newToken.GetId(),
			).Execute()
			assert.NoError(t, err)
			assert.Equal(t, "disabled", *checkProp.GetProperties().GetStatus())
			time.Sleep(20 * time.Second)
			// delete token
			dToken := TokenDeleteCmd()
			dToken.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			dToken.Command.Flags().Set(FlagTokenId, newToken.GetId())

			err = dToken.Command.Execute()
			assert.NoError(t, err)
			time.Sleep(20 * time.Second)

			// replace token
			rToken := TokenReplaceCmd()
			rToken.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())
			rToken.Command.Flags().Set(FlagTokenId, newToken.GetId())
			rToken.Command.Flags().Set(constants.FlagName, "newName")
			time.Sleep(20 * time.Second)

			err = rToken.Command.Execute()
			assert.NoError(t, err)
			time.Sleep(20 * time.Second)

			// delete registry
			d := registry.RegDeleteCmd()
			d.Command.Flags().Set(constants.FlagRegistryId, newReg.GetId())

			err = d.Command.Execute()
			assert.NoError(t, err)

		},
	)
}
