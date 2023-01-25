package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	"github.com/fatih/structs"
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
	List(params ListQueryParams) (BackupUnits, *Response, error)
	Get(backupUnitId string, params QueryParams) (*BackupUnit, *Response, error)
	GetSsoUrl(backupUnitId string, params QueryParams) (*BackupUnitSSO, *Response, error)
	Create(u BackupUnit, params QueryParams) (*BackupUnit, *Response, error)
	Update(backupUnitId string, input BackupUnitProperties, params QueryParams) (*BackupUnit, *Response, error)
	Delete(backupUnitId string, params QueryParams) (*Response, error)
}

type backupUnitsService struct {
	client  *config.Client
	context context.Context
}

var _ BackupUnitsService = &backupUnitsService{}

func NewBackupUnitService(client *config.Client, ctx context.Context) BackupUnitsService {
	return &backupUnitsService{
		client:  client,
		context: ctx,
	}
}

func (s *backupUnitsService) List(params ListQueryParams) (BackupUnits, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				req = req.Filter(k, v)
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	dcs, res, err := s.client.BackupUnitsApi.BackupunitsGetExecute(req)
	return BackupUnits{dcs}, &Response{*res}, err
}

func (s *backupUnitsService) Get(backupUnitId string, params QueryParams) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsFindById(s.context, backupUnitId)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsFindByIdExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) GetSsoUrl(backupUnitId string, params QueryParams) (*BackupUnitSSO, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsSsourlGet(s.context, backupUnitId)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsSsourlGetExecute(req)
	return &BackupUnitSSO{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Create(u BackupUnit, params QueryParams) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsPost(s.context).BackupUnit(u.BackupUnit)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsPostExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Update(backupUnitId string, input BackupUnitProperties, params QueryParams) (*BackupUnit, *Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsPatch(s.context, backupUnitId).BackupUnit(input.BackupUnitProperties)
	backupUnit, res, err := s.client.BackupUnitsApi.BackupunitsPatchExecute(req)
	return &BackupUnit{backupUnit}, &Response{*res}, err
}

func (s *backupUnitsService) Delete(backupUnitId string, params QueryParams) (*Response, error) {
	req := s.client.BackupUnitsApi.BackupunitsDelete(s.context, backupUnitId)
	res, err := s.client.BackupUnitsApi.BackupunitsDeleteExecute(req)
	return &Response{*res}, err
}
