package dbaas_postgres

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	mockResources "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources/mocks"
)

type ResourcesMocks struct {
	Cluster *mockResources.MockClustersService
	Version *mockResources.MockVersionsService
	Info    *mockResources.MockInfosService
	Log     *mockResources.MockLogsService
	Restore *mockResources.MockRestoresService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Cluster: mockResources.NewMockClustersService(ctrl),
		Version: mockResources.NewMockVersionsService(ctrl),
		Info:    mockResources.NewMockInfosService(ctrl),
		Log:     mockResources.NewMockLogsService(ctrl),
		Restore: mockResources.NewMockRestoresService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Clusters = func() resources.ClustersService { return tm.Cluster }
	c.Versions = func() resources.VersionsService { return tm.Version }
	c.Infos = func() resources.InfosService { return tm.Info }
	c.Logs = func() resources.LogsService { return tm.Log }
	c.Restores = func() resources.RestoresService { return tm.Restore }
	return c
}
