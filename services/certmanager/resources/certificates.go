package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

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
	Delete(certId string) (*sdkgo.APIResponse, error)
	Patch(certId string, input sdkgo.CertificatePatchDto) (sdkgo.CertificateDto, *sdkgo.APIResponse, error)
	GetApiVersion() (sdkgo.ApiInfoDto, *sdkgo.APIResponse, error)
}

type certsService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ CertsService = &certsService{}

func NewCertsService(client *config.Client, ctx context.Context) CertsService {
	return &certsService{
		client:  client.CertManagerClient,
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

func (svc *certsService) Delete(certId string) (*sdkgo.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesDelete(svc.context, certId)
	res, err := svc.client.CertificatesApi.CertificatesDeleteExecute(req)
	return res, err
}

func (svc *certsService) Patch(certId string, input sdkgo.CertificatePatchDto) (sdkgo.CertificateDto, *sdkgo.APIResponse, error) {
	req := svc.client.CertificatesApi.CertificatesPatch(svc.context, certId).CertificatePatchDto(input)
	cert, res, err := svc.client.CertificatesApi.CertificatesPatchExecute(req)
	return cert, res, err
}

func (svc *certsService) GetApiVersion() (sdkgo.ApiInfoDto, *sdkgo.APIResponse, error) {
	req := svc.client.InformationApi.GetInfo(svc.context)
	api, res, err := svc.client.InformationApi.GetInfoExecute(req)
	return api, res, err
}
