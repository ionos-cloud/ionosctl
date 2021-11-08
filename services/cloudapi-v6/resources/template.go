package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Template struct {
	ionoscloud.Template
}

type TemplateProperties struct {
	ionoscloud.TemplateProperties
}

type Templates struct {
	ionoscloud.Templates
}

// TemplatesService is a wrapper around ionoscloud.Template
type TemplatesService interface {
	List(params ListQueryParams) (Templates, *Response, error)
	Get(templateId string) (*Template, *Response, error)
}

type templatesService struct {
	client  *Client
	context context.Context
}

var _ TemplatesService = &templatesService{}

func NewTemplateService(client *Client, ctx context.Context) TemplatesService {
	return &templatesService{
		client:  client,
		context: ctx,
	}
}

func (ss *templatesService) List(params ListQueryParams) (Templates, *Response, error) {
	req := ss.client.TemplatesApi.TemplatesGet(ss.context)
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
	s, res, err := ss.client.TemplatesApi.TemplatesGetExecute(req)
	return Templates{s}, &Response{*res}, err
}

func (ss *templatesService) Get(templateId string) (*Template, *Response, error) {
	req := ss.client.TemplatesApi.TemplatesFindById(ss.context, templateId)
	s, res, err := ss.client.TemplatesApi.TemplatesFindByIdExecute(req)
	return &Template{s}, &Response{*res}, err
}
