package v5

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	client  *Client
	context context.Context
}

var _ BackupUnitsService = &backupUnitsService{}

func NewBackupUnitService(client *Client, ctx context.Context) BackupUnitsService {
	return &backupUnitsService{
		client:  client,
		context: ctx,
	}
}

func (s *backupUnitsService) List() (BackupUnits, *Response, error) {
	req := s.client.BackupUnitApi.BackupunitsGet(s.context)
	dcs, res, err := s.client.BackupUnitApi.BackupunitsGetExecute(req)
	return BackupUnits{dcs}, &Response{*res}, err
}

func (s *backupUnitsService) Get(backupUnitId string) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitApi.BackupunitsFindById(s.context, backupUnitId)
	backupUnit, res, err := s.client.BackupUnitApi.BackupunitsFindByIdExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) GetSsoUrl(backupUnitId string) (*BackupUnitSSO, *Response, error) {
	req := s.client.BackupUnitApi.BackupunitsSsourlGet(s.context, backupUnitId)
	backupUnit, res, err := s.client.BackupUnitApi.BackupunitsSsourlGetExecute(req)
	return &BackupUnitSSO{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Create(u BackupUnit) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitApi.BackupunitsPost(s.context).BackupUnit(u.BackupUnit)
	backupUnit, res, err := s.client.BackupUnitApi.BackupunitsPostExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Update(backupUnitId string, input BackupUnitProperties) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitApi.BackupunitsPatch(s.context, backupUnitId).BackupUnitProperties(input.BackupUnitProperties)
	backupUnit, res, err := s.client.BackupUnitApi.BackupunitsPatchExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Delete(backupUnitId string) (*Response, error) {
	req := s.client.BackupUnitApi.BackupunitsDelete(s.context, backupUnitId)
	_, res, err := s.client.BackupUnitApi.BackupunitsDeleteExecute(req)
	return &Response{*res}, err
}
