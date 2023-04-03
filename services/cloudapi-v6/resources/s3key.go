package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
)

type S3Key struct {
	ionoscloud.S3Key
}

type S3Keys struct {
	ionoscloud.S3Keys
}

// S3KeysService is a wrapper around ionoscloud.S3Key
type S3KeysService interface {
	List(userId string, params ListQueryParams) (S3Keys, *Response, error)
	Get(userId, keyId string, params QueryParams) (*S3Key, *Response, error)
	Create(userId string, params QueryParams) (*S3Key, *Response, error)
	Update(userId, keyId string, key S3Key, params QueryParams) (*S3Key, *Response, error)
	Delete(userId, keyId string, params QueryParams) (*Response, error)
}

type s3KeysService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ S3KeysService = &s3KeysService{}

func NewS3KeyService(client *client.Client, ctx context.Context) S3KeysService {
	return &s3KeysService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *s3KeysService) List(userId string, params ListQueryParams) (S3Keys, *Response, error) {
	req := s.client.UserS3KeysApi.UmUsersS3keysGet(s.context, userId)
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
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	keys, resp, err := s.client.UserS3KeysApi.UmUsersS3keysGetExecute(req)
	return S3Keys{keys}, &Response{*resp}, err
}

func (s *s3KeysService) Get(userId, keyId string, params QueryParams) (*S3Key, *Response, error) {
	req := s.client.UserS3KeysApi.UmUsersS3keysFindByKeyId(s.context, userId, keyId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	key, resp, err := s.client.UserS3KeysApi.UmUsersS3keysFindByKeyIdExecute(req)
	return &S3Key{key}, &Response{*resp}, err
}

func (s *s3KeysService) Create(userId string, params QueryParams) (*S3Key, *Response, error) {
	req := s.client.UserS3KeysApi.UmUsersS3keysPost(s.context, userId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s3key, resp, err := s.client.UserS3KeysApi.UmUsersS3keysPostExecute(req)
	return &S3Key{s3key}, &Response{*resp}, err
}

func (s *s3KeysService) Update(userId, keyId string, key S3Key, params QueryParams) (*S3Key, *Response, error) {
	req := s.client.UserS3KeysApi.UmUsersS3keysPut(s.context, userId, keyId).S3Key(key.S3Key)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s3key, resp, err := s.client.UserS3KeysApi.UmUsersS3keysPutExecute(req)
	return &S3Key{s3key}, &Response{*resp}, err
}

func (s *s3KeysService) Delete(userId, keyId string, params QueryParams) (*Response, error) {
	req := s.client.UserS3KeysApi.UmUsersS3keysDelete(s.context, userId, keyId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	resp, err := s.client.UserS3KeysApi.UmUsersS3keysDeleteExecute(req)
	return &Response{*resp}, err
}
