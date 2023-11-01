//go:build unit
// +build unit

package container_registry

import (
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/location"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/name"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestContainerRegistryServiceCmd(t *testing.T) {
	core.RootCmdTest.AddCommand(ContainerRegistryCmd())
	assert.True(t, ContainerRegistryCmd().IsAvailableCommand())
	assert.True(t, registry.RegPostCmd().IsAvailableCommand())
	assert.True(t, registry.RegGetCmd().IsAvailableCommand())
	assert.True(t, registry.RegListCmd().IsAvailableCommand())
	assert.True(t, registry.RegDeleteCmd().IsAvailableCommand())
	assert.True(t, registry.RegUpdateCmd().IsAvailableCommand())
	assert.True(t, registry.RegReplaceCmd().IsAvailableCommand())
	assert.True(t, token.TokenPostCmd().IsAvailableCommand())
	assert.True(t, token.TokenGetCmd().IsAvailableCommand())
	assert.True(t, token.TokenListCmd().IsAvailableCommand())
	assert.True(t, token.TokenDeleteCmd().IsAvailableCommand())
	assert.True(t, token.TokenUpdateCmd().IsAvailableCommand())
	assert.True(t, token.TokenReplaceCmd().IsAvailableCommand())
	assert.True(t, name.RegNamesCmd().IsAvailableCommand())
	assert.True(t, location.RegLocationsListCmd().IsAvailableCommand())
	assert.True(t, repository.RegRepoDeleteCmd().IsAvailableCommand())
}
