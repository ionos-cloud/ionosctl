package dataplatform

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"
	mockResources "github.com/ionos-cloud/ionosctl/services/dataplatform/resources/mocks"
)

type ResourcesMocks struct {
	Client  *mockResources.MockClientService
	Cluster *mockResources.MockClustersService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Client:  mockResources.NewMockClientService(ctrl),
		Cluster: mockResources.NewMockClustersService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Clusters = func() resources.ClustersService { return tm.Cluster }
	return c
}
