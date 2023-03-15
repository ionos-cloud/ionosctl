package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type TemplatesService interface {
	List(limit, offset *int32) (sdkgo.TemplateList, *sdkgo.APIResponse, error)
}

type templatesService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ TemplatesService = &templatesService{}

func NewTemplatesService(client *client.Client, ctx context.Context) TemplatesService {
	return &templatesService{
		client:  client.MongoClient,
		context: ctx,
	}
}

func (svc *templatesService) List(limit, offset *int32) (sdkgo.TemplateList, *sdkgo.APIResponse, error) {
	req := svc.client.TemplatesApi.TemplatesGet(svc.context)
	if limit != nil {
		req = req.Limit(*limit)
	}
	if offset != nil {
		req = req.Offset(*offset)
	}
	ls, res, err := svc.client.TemplatesApi.TemplatesGetExecute(req)
	return ls, res, err
}
