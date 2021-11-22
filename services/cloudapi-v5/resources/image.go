package resources

import (
	"context"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Image struct {
	ionoscloud.Image
}

type Images struct {
	ionoscloud.Images
}

type ImageProperties struct {
	ionoscloud.ImageProperties
}

// ImagesService is a wrapper around ionoscloud.Image
type ImagesService interface {
	List(params ListQueryParams) (Images, *Response, error)
	Get(imageId string) (*Image, *Response, error)
	Update(imageId string, imgProp ImageProperties) (*Image, *Response, error)
	Delete(imageId string) (*Response, error)
}

type imagesService struct {
	client  *Client
	context context.Context
}

var _ ImagesService = &imagesService{}

func NewImageService(client *Client, ctx context.Context) ImagesService {
	return &imagesService{
		client:  client,
		context: ctx,
	}
}

func (s *imagesService) List(params ListQueryParams) (Images, *Response, error) {
	req := s.client.ImageApi.ImagesGet(s.context)
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
	}
	images, resp, err := s.client.ImageApi.ImagesGetExecute(req)
	return Images{images}, &Response{*resp}, err
}

func (s *imagesService) Get(imageId string) (*Image, *Response, error) {
	req := s.client.ImageApi.ImagesFindById(s.context, imageId)
	image, resp, err := s.client.ImageApi.ImagesFindByIdExecute(req)
	return &Image{image}, &Response{*resp}, err
}

func (s *imagesService) Update(imageId string, imgProp ImageProperties) (*Image, *Response, error) {
	req := s.client.ImageApi.ImagesPatch(s.context, imageId).Image(imgProp.ImageProperties)
	image, resp, err := s.client.ImageApi.ImagesPatchExecute(req)
	return &Image{image}, &Response{*resp}, err
}

func (s *imagesService) Delete(imageId string) (*Response, error) {
	req := s.client.ImageApi.ImagesDelete(s.context, imageId)
	_, resp, err := s.client.ImageApi.ImagesDeleteExecute(req)
	return &Response{*resp}, err
}
