package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type TemplatesService interface {
	List() (sdkgo.TemplateList, *sdkgo.APIResponse, error)
}

type templatesService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ TemplatesService = &templatesService{}

func NewTemplatesService(client *config.Client, ctx context.Context) TemplatesService {
	return &templatesService{
		client:  client.MongoClient,
		context: ctx,
	}
}

func (svc *templatesService) List() (sdkgo.TemplateList, *sdkgo.APIResponse, error) {
	req := svc.client.TemplatesApi.TemplatesGet(svc.context)
	ls, res, err := svc.client.TemplatesApi.TemplatesGetExecute(req)
	return ls, res, err
}
