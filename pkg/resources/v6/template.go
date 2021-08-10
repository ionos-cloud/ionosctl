package v6

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
	List() (Templates, *Response, error)
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

func (ss *templatesService) List() (Templates, *Response, error) {
	req := ss.client.TemplatesApi.TemplatesGet(ss.context)
	s, res, err := ss.client.TemplatesApi.TemplatesGetExecute(req)
	return Templates{s}, &Response{*res}, err
}

func (ss *templatesService) Get(templateId string) (*Template, *Response, error) {
	req := ss.client.TemplatesApi.TemplatesFindById(ss.context, templateId)
	s, res, err := ss.client.TemplatesApi.TemplatesFindByIdExecute(req)
	return &Template{s}, &Response{*res}, err
}
