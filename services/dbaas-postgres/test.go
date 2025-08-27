package dbaas_postgres

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	mockResources "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources/mocks"
)

type ResourcesMocks struct {
	Cluster *mockResources.MockClustersService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Cluster: mockResources.NewMockClustersService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Clusters = func() resources.ClustersService { return tm.Cluster }
	return c
}
