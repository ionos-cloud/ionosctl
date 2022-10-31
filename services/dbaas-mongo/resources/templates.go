package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type TemplatesService interface {
	List() (sdkgo.TemplateList, *sdkgo.APIResponse, error)
}

type templatesService struct {
	client  *Client
	context context.Context
}

var _ TemplatesService = &templatesService{}

func NewTemplatesService(client *Client, ctx context.Context) TemplatesService {
	return &templatesService{
		client:  client,
		context: ctx,
	}
}

func (svc *templatesService) List() (sdkgo.TemplateList, *sdkgo.APIResponse, error) {
	req := svc.client.TemplatesApi.TemplatesGet(svc.context)
	ls, res, err := svc.client.TemplatesApi.TemplatesGetExecute(req)
	return ls, res, err
}
