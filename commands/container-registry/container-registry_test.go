package container_registry

import (
	"errors"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/location"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/name"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/token"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainerRegistryServiceCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ContainerRegistryCmd())
	if ok := ContainerRegistryCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := registry.RegPostCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := registry.RegGetCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := registry.RegListCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := registry.RegDeleteCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := registry.RegUpdateCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := registry.RegReplaceCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := token.TokenPostCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := token.TokenGetCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := token.TokenListCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := token.TokenDeleteCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := token.TokenUpdateCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := token.TokenReplaceCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := name.RegNamesCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := location.RegLocationsListCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
	if ok := repository.RegRepoDeleteCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
}
