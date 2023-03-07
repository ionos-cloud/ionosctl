package auth_v1

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
	mockResources "github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources/mocks"
)

type ResourcesMocks struct {
	Client *mockResources.MockClientService
	Token  *mockResources.MockTokensService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Client: mockResources.NewMockClientService(ctrl),
		Token:  mockResources.NewMockTokensService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Tokens = func() resources.TokensService { return tm.Token }
	return c
}
