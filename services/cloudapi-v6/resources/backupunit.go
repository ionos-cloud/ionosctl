package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type BackupUnit struct {
	ionoscloud.BackupUnit
}

type BackupUnitSSO struct {
	ionoscloud.BackupUnitSSO
}

type BackupUnitProperties struct {
	ionoscloud.BackupUnitProperties
}

type BackupUnits struct {
	ionoscloud.BackupUnits
}

// BackupUnitsService is a wrapper around ionoscloud.BackupUnit
type BackupUnitsService interface {
	List() (BackupUnits, *Response, error)
	Get(backupUnitId string) (*BackupUnit, *Response, error)
	GetSsoUrl(backupUnitId string) (*BackupUnitSSO, *Response, error)
	Create(u BackupUnit) (*BackupUnit, *Response, error)
	Update(backupUnitId string, input BackupUnitProperties) (*BackupUnit, *Response, error)
	Delete(backupUnitId string) (*Response, error)
}

type backupUnitsService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ BackupUnitsService = &backupUnitsService{}

func NewBackupUnitService(client *client.Client, ctx context.Context) BackupUnitsService {
	return &backupUnitsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *backupUnitsService) List() (BackupUnits, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsGet(s.context)
	dcs, res, err := s.client.BackupUnitsApi.BackupunitsGetExecute(req)
	return BackupUnits{dcs}, &Response{*res}, err
}

func (s *backupUnitsService) Get(backupUnitId string) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsFindById(s.context, backupUnitId)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsFindByIdExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) GetSsoUrl(backupUnitId string) (*BackupUnitSSO, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsSsourlGet(s.context, backupUnitId)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsSsourlGetExecute(req)
	return &BackupUnitSSO{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Create(u BackupUnit) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsPost(s.context).BackupUnit(u.BackupUnit)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsPostExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Update(backupUnitId string, input BackupUnitProperties) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsPatch(s.context, backupUnitId).BackupUnit(input.BackupUnitProperties)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsPatchExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Delete(backupUnitId string) (*Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsDelete(s.context, backupUnitId)
	res, err := s.client.BackupUnitsApi.BackupunitsDeleteExecute(req)
	return &Response{*res}, err
}
