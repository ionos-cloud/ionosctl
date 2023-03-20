//go:build integration
// +build integration

package registry

import (
	"context"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/lucasjones/reggen"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
		c := RegPostCmd()
		c.Command.Flags().Set(FlagName, name)
		c.Command.Flags().Set(FlagLocation, "de/fra")

		err = c.Command.Execute()
		assert.NoError(t, err)

		registries, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
		assert.NoError(t, err)

		var newRegistry *ionoscloud.RegistryResponse
		for _, registry := range *registries.GetItems() {
			if *registry.GetProperties().GetName() == name {
				newRegistry = &registry
			}
		}

		g := RegGetCmd()
		g.Command.Flags().Set(FlagRegId, *newRegistry.GetId())
		assert.NoError(t, err)

		err = g.Command.Execute()
		assert.NoError(t, err)

		patch := RegUpdateCmd()
		patch.Command.Flags().Set(FlagRegId, *newRegistry.GetId())
		patch.Command.Flags().Set(FlagRegGCDays, "Tuesday")
		assert.NoError(t, err)

		err = patch.Command.Execute()
		assert.NoError(t, err)

		checkRegistry, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesFindById(context.Background(), *newRegistry.GetId()).Execute()
		assert.NoError(t, err)
		assert.Equal(t, []ionoscloud.Day([]ionoscloud.Day{"Tuesday"}), *checkRegistry.GetProperties().GetGarbageCollectionSchedule().GetDays())

		d := RegDeleteCmd()
		d.Command.Flags().Set(FlagRegId, *newRegistry.GetId())
		assert.NoError(t, err)

		err = d.Command.Execute()
		assert.NoError(t, err)

		replace := RegReplaceCmd()
		name, err = reggen.Generate("^[a-z][-a-z0-9]{1,61}[a-z0-9]$", 10)
		assert.NoError(t, err)
		replace.Command.Flags().Set(FlagRegId, *newRegistry.GetId())
		replace.Command.Flags().Set(FlagName, name)
		replace.Command.Flags().Set(FlagLocation, "de/fra")
		replace.Command.Flags().Set(FlagRegGCDays, "Tuesday")

		err = replace.Command.Execute()
		assert.NoError(t, err)

		d.Command.Flags().Set(FlagRegId, *newRegistry.GetId())

		err = d.Command.Execute()
		assert.NoError(t, err)
	})
}
