package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	cert "github.com/ionos-cloud/sdk-go-cert-manager"
)

type Response struct {
	cert.APIResponse
}

// CertsService is a wrapper around ionoscloud.CertificateDto
type CertsService interface {
	Get(certId string) (cert.CertificateDto, *cert.APIResponse, error)
	Post(cert.CertificatePostDto) (cert.CertificateDto, *cert.APIResponse, error)
	List() (cert.CertificateCollectionDto, *cert.APIResponse, error)
	Delete(certId string) (*cert.APIResponse, error)
	Patch(certId string, input cert.CertificatePatchDto) (cert.CertificateDto, *cert.APIResponse, error)
	GetApiVersion() (cert.ApiInfoDto, *cert.APIResponse, error)
}

type certsService struct {
	client  *cert.APIClient
	context context.Context
}

var _ CertsService = &certsService{}

func NewCertsService(client *client.Client, ctx context.Context) CertsService {
	return &certsService{
		client:  client.CertManagerClient,
		context: ctx,
	}
}

func (svc *certsService) Get(certId string) (cert.CertificateDto, *cert.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesGetById(svc.context, certId)
	cert, res, err := svc.client.CertificatesApi.CertificatesGetByIdExecute(req)
	return cert, res, err
}

func (svc *certsService) Post(input cert.CertificatePostDto) (cert.CertificateDto, *cert.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesPost(svc.context).CertificatePostDto(input)
	cert, res, err := svc.client.CertificatesApi.CertificatesPostExecute(req)
	return cert, res, err
}

func (svc *certsService) List() (cert.CertificateCollectionDto, *cert.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesGet(svc.context)
	cert, res, err := svc.client.CertificatesApi.CertificatesGetExecute(req)
	return cert, res, err
}

func (svc *certsService) Delete(certId string) (*cert.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesDelete(svc.context, certId)
	res, err := svc.client.CertificatesApi.CertificatesDeleteExecute(req)
	return res, err
}

func (svc *certsService) Patch(certId string, input cert.CertificatePatchDto) (cert.CertificateDto, *cert.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesPatch(svc.context, certId).CertificatePatchDto(input)
	cert, res, err := svc.client.CertificatesApi.CertificatesPatchExecute(req)
	return cert, res, err
}

func (svc *certsService) GetApiVersion() (cert.ApiInfoDto, *cert.APIResponse, error) {
	req := svc.client.InformationApi.GetInfo(svc.context)
	api, res, err := svc.client.InformationApi.GetInfoExecute(req)
	return api, res, err
}
