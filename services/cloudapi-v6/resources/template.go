package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
)

type Template struct {
	compute.Template
}

type TemplateProperties struct {
	compute.TemplateProperties
}

type Templates struct {
	compute.Templates
}

// TemplatesService is a wrapper around compute.Template
type TemplatesService interface {
	List(params ListQueryParams) (Templates, *Response, error)
	Get(templateId string, params QueryParams) (*Template, *Response, error)
}

type templatesService struct {
	client  *compute.APIClient
	context context.Context
}

var _ TemplatesService = &templatesService{}

func NewTemplateService(client *client.Client, ctx context.Context) TemplatesService {
	return &templatesService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (ss *templatesService) List(params ListQueryParams) (Templates, *Response, error) {
	req := ss.client.TemplatesApi.TemplatesGet(ss.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
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
			// if params.QueryParams.Pretty != nil {
			//	// Currently not implemented
			//	req = req.Pretty(*params.QueryParams.Pretty)
			// }
		}
	}
	s, res, err := ss.client.TemplatesApi.TemplatesGetExecute(req)
	return Templates{s}, &Response{*res}, err
}

func (ss *templatesService) Get(templateId string, params QueryParams) (*Template, *Response, error) {
	req := ss.client.TemplatesApi.TemplatesFindById(ss.context, templateId)
	s, res, err := ss.client.TemplatesApi.TemplatesFindByIdExecute(req)
	return &Template{s}, &Response{*res}, err
}
