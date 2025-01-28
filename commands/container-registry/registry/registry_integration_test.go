//go:build integration
// +build integration

package registry

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func teardown() {
	regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()

	if err != nil {
		log.Print(fmt.Errorf("failed deleting all registries: %w", err))
	}
	for _, reg := range *regs.Items {
		_, err := client.Must().RegistryClient.RegistriesApi.RegistriesDelete(context.Background(), *reg.Id).Execute()
		if err != nil {
			log.Print(fmt.Errorf("failed deleting registry: %w", err))
		}
	}
	time.Sleep(30 * time.Second)
}

func TestRegistryService(t *testing.T) {
	t.Run(
		"registry functions", func(t *testing.T) {
			t.Cleanup(teardown)
			viper.Reset()
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgForce, true)

			name := "ionosctl-crreg-test-" + fake.AlphaNum(8)
			c := RegPostCmd()
			c.Command.Flags().Set(constants.FlagName, name)
			c.Command.Flags().Set(constants.FlagLocation, "de/fra")

			err := c.Command.Execute()
			assert.NoError(t, err)

			registries, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
			assert.NoError(t, err)

			var newRegistry *containerregistry.RegistryResponse
			for _, registry := range *registries.GetItems() {
				if *registry.GetProperties().GetName() == name {
					newRegistry = &registry
				}
			}

			g := RegGetCmd()
			g.Command.Flags().Set(constants.FlagRegistryId, *newRegistry.GetId())
			assert.NoError(t, err)

			err = g.Command.Execute()
			assert.NoError(t, err)

			patch := RegUpdateCmd()
			patch.Command.Flags().Set(constants.FlagRegistryId, *newRegistry.GetId())
			patch.Command.Flags().Set(FlagRegGCDays, "Tuesday")
			assert.NoError(t, err)

			err = patch.Command.Execute()
			assert.NoError(t, err)

			checkRegistry, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesFindById(
				context.Background(), *newRegistry.GetId(),
			).Execute()
			assert.NoError(t, err)
			assert.Equal(
				t, []containerregistry.Day([]containerregistry.Day{"Tuesday"}),
				*checkRegistry.GetProperties().GetGarbageCollectionSchedule().GetDays(),
			)

			d := RegDeleteCmd()
			d.Command.Flags().Set(constants.FlagRegistryId, *newRegistry.GetId())
			assert.NoError(t, err)

			err = d.Command.Execute()
			assert.NoError(t, err)

			replace := RegReplaceCmd()
			name = "ionosctl-crreg-test-" + fake.AlphaNum(8)
			replace.Command.Flags().Set(constants.FlagRegistryId, *newRegistry.GetId())
			replace.Command.Flags().Set(constants.FlagName, name)
			replace.Command.Flags().Set(constants.FlagLocation, "de/fra")
			replace.Command.Flags().Set(FlagRegGCDays, "Tuesday")

			err = replace.Command.Execute()
			assert.NoError(t, err)

			d.Command.Flags().Set(constants.FlagRegistryId, *newRegistry.GetId())

			err = d.Command.Execute()
			assert.NoError(t, err)
		},
	)
}
