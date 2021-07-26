package v5

import (
	"context"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type S3Key struct {
	ionoscloud.S3Key
}

type S3Keys struct {
	ionoscloud.S3Keys
}

// S3KeysService is a wrapper around ionoscloud.S3Key
type S3KeysService interface {
	List(userId string) (S3Keys, *Response, error)
	Get(userId, keyId string) (*S3Key, *Response, error)
	Create(userId string) (*S3Key, *Response, error)
	Update(userId, keyId string, key S3Key) (*S3Key, *Response, error)
	Delete(userId, keyId string) (*Response, error)
}

type s3KeysService struct {
	client  *Client
	context context.Context
}

var _ S3KeysService = &s3KeysService{}

func NewS3KeyService(client *Client, ctx context.Context) S3KeysService {
	return &s3KeysService{
		client:  client,
		context: ctx,
	}
}

func (s *s3KeysService) List(userId string) (S3Keys, *Response, error) {
	req := s.client.UserManagementApi.UmUsersS3keysGet(s.context, userId)
	keys, resp, err := s.client.UserManagementApi.UmUsersS3keysGetExecute(req)
	return S3Keys{keys}, &Response{*resp}, err
}

func (s *s3KeysService) Get(userId, keyId string) (*S3Key, *Response, error) {
	req := s.client.UserManagementApi.UmUsersS3keysFindByKeyId(s.context, userId, keyId)
	key, resp, err := s.client.UserManagementApi.UmUsersS3keysFindByKeyIdExecute(req)
	return &S3Key{key}, &Response{*resp}, err
}

func (s *s3KeysService) Create(userId string) (*S3Key, *Response, error) {
	req := s.client.UserManagementApi.UmUsersS3keysPost(s.context, userId)
	s3key, resp, err := s.client.UserManagementApi.UmUsersS3keysPostExecute(req)
	return &S3Key{s3key}, &Response{*resp}, err
}

func (s *s3KeysService) Update(userId, keyId string, key S3Key) (*S3Key, *Response, error) {
	req := s.client.UserManagementApi.UmUsersS3keysPut(s.context, userId, keyId).S3Key(key.S3Key)
	s3key, resp, err := s.client.UserManagementApi.UmUsersS3keysPutExecute(req)
	return &S3Key{s3key}, &Response{*resp}, err
}

func (s *s3KeysService) Delete(userId, keyId string) (*Response, error) {
	req := s.client.UserManagementApi.UmUsersS3keysDelete(s.context, userId, keyId)
	_, resp, err := s.client.UserManagementApi.UmUsersS3keysDeleteExecute(req)
	return &Response{*resp}, err
}
