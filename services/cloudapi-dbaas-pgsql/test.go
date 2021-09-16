package cloudapi_dbaas_pgsql

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
	mockResources "github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources/mocks"
)

type ResourcesMocks struct {
	Client  *mockResources.MockClientService
	Cluster *mockResources.MockClustersService
	Backup  *mockResources.MockBackupsService
	Version *mockResources.MockVersionsService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Client:  mockResources.NewMockClientService(ctrl),
		Cluster: mockResources.NewMockClustersService(ctrl),
		Backup:  mockResources.NewMockBackupsService(ctrl),
		Version: mockResources.NewMockVersionsService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Clusters = func() resources.ClustersService { return tm.Cluster }
	c.Backups = func() resources.BackupsService { return tm.Backup }
	c.Versions = func() resources.VersionsService { return tm.Version }
	return c
}
