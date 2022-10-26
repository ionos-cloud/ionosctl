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
	GetById(certId string) (sdkgo.CertificateDto, *sdkgo.APIResponse, error)
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

func (svc *certsService) GetById(certId string) (sdkgo.CertificateDto, *sdkgo.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesGetById(svc.context, certId)
	cert, res, err := svc.client.CertificatesApi.CertificatesGetByIdExecute(req)
	return cert, res, err
}
