package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
)

type Response struct {
	sdkgo.APIResponse
}

// CertsService is a wrapper around ionoscloud.CertificateDto
type CertsService interface {
	Get(certId string) (sdkgo.CertificateDto, *sdkgo.APIResponse, error)
	Post(sdkgo.CertificatePostDto) (sdkgo.CertificateDto, *sdkgo.APIResponse, error)
	List() (sdkgo.CertificateCollectionDto, *sdkgo.APIResponse, error)
}

type certsService struct {
	client  *Client
	context context.Context
}

var _ CertsService = &certsService{}

func NewCertsService(client *Client, ctx context.Context) CertsService {
	return &certsService{
		client:  client,
		context: ctx,
	}
}

func (svc *certsService) Get(certId string) (sdkgo.CertificateDto, *sdkgo.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesGetById(svc.context, certId)
	cert, res, err := svc.client.CertificatesApi.CertificatesGetByIdExecute(req)
	return cert, res, err
}

func (svc *certsService) Post(input sdkgo.CertificatePostDto) (sdkgo.CertificateDto, *sdkgo.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesPost(svc.context).CertificatePostDto(input)
	cert, res, err := svc.client.CertificatesApi.CertificatesPostExecute(req)
	return cert, res, err
}

func (svc *certsService) List() (sdkgo.CertificateCollectionDto, *sdkgo.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesGet(svc.context)
	cert, res, err := svc.client.CertificatesApi.CertificatesGetExecute(req)
	return cert, res, err
}
