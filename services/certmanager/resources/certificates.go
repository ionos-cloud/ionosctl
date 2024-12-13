package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
)

type Response struct {
	shared.APIResponse
}

// CertsService is a wrapper around ionoscloud.Certificate
type CertsService interface {
	Get(certId string) (sdkgo.Certificate, *shared.APIResponse, error)
	Post(sdkgo.CertificateCreate) (sdkgo.Certificate, *shared.APIResponse, error)
	List() (sdkgo.CertificateReadList, *shared.APIResponse, error)
	Delete(certId string) (*shared.APIResponse, error)
	Patch(certId string, input sdkgo.CertificatePatch) (sdkgo.Certificate, *shared.APIResponse, error)
	// GetApiVersion() (sdkgo.ApiInfoDto, *shared.APIResponse, error)
}

type certsService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ CertsService = &certsService{}

func NewCertsService(client *client.Client, ctx context.Context) CertsService {
	return &certsService{
		client:  client.CertManagerClient,
		context: ctx,
	}
}

func (svc *certsService) Get(certId string) (sdkgo.Certificate, *shared.APIResponse, error) {
	req := svc.client.CertificateApi.CertificatesFindById(svc.context, certId)
	cert, res, err := svc.client.CertificateApi.CertificatesFindByIdExecute(req)
	return cert, res, err
}

func (svc *certsService) Post(input sdkgo.CertificateCreate) (sdkgo.Certificate, *shared.APIResponse, error) {
	req := svc.client.CertificateApi.CertificatesPost(svc.context).CertificateCreate(input)
	cert, res, err := svc.client.CertificateApi.CertificatesPostExecute(req)
	return cert, res, err
}

func (svc *certsService) List() (sdkgo.CertificateReadList, *shared.APIResponse, error) {
	req := svc.client.CertificateApi.CertificatesGet(svc.context)
	cert, res, err := svc.client.CertificateApi.CertificatesGetExecute(req)
	return cert, res, err
}

func (svc *certsService) Delete(certId string) (*shared.APIResponse, error) {
	req := svc.client.CertificateApi.CertificatesDelete(svc.context, certId)
	res, err := svc.client.CertificateApi.CertificatesDeleteExecute(req)
	return res, err
}

func (svc *certsService) Patch(certId string, input sdkgo.CertificatePatch) (sdkgo.Certificate, *shared.APIResponse, error) {
	req := svc.client.CertificateApi.CertificatesPatch(svc.context, certId).CertificatePatch(input)
	cert, res, err := svc.client.CertificateApi.CertificatesPatchExecute(req)
	return cert, res, err
}

func (svc *certsService) GetApiVersion() (sdkgo.Metadata, *shared.APIResponse, error) {
	req := svc.client.ProviderApi.(svc.context)
	api, res, err := svc.client.InformationApi.GetInfoExecute(req)
	return api, res, err
}
