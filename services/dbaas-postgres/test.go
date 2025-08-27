package dbaas_postgres

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	mockResources "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources/mocks"
)

type ResourcesMocks struct {
	Cluster *mockResources.MockClustersService
	Log     *mockResources.MockLogsService
	Restore *mockResources.MockRestoresService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Cluster: mockResources.NewMockClustersService(ctrl),
		Log:     mockResources.NewMockLogsService(ctrl),
		Restore: mockResources.NewMockRestoresService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Clusters = func() resources.ClustersService { return tm.Cluster }
	c.Logs = func() resources.LogsService { return tm.Log }
	c.Restores = func() resources.RestoresService { return tm.Restore }
	return c
}
