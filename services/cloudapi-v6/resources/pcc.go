package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
)

type PrivateCrossConnect struct {
	ionoscloud.PrivateCrossConnect
}

type PrivateCrossConnectProperties struct {
	ionoscloud.PrivateCrossConnectProperties
}

type PrivateCrossConnects struct {
	ionoscloud.PrivateCrossConnects
}

type Peer struct {
	ionoscloud.Peer
}

// PccsService is a wrapper around ionoscloud.PrivateCrossConnect
type PccsService interface {
	List(params ListQueryParams) (PrivateCrossConnects, *Response, error)
	Get(pccId string, params QueryParams) (*PrivateCrossConnect, *Response, error)
	GetPeers(pccId string) (*[]Peer, *Response, error)
	Create(u PrivateCrossConnect, params QueryParams) (*PrivateCrossConnect, *Response, error)
	Update(pccId string, input PrivateCrossConnectProperties, params QueryParams) (*PrivateCrossConnect, *Response, error)
	Delete(pccId string, params QueryParams) (*Response, error)
}

type pccsService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ PccsService = &pccsService{}

func NewPrivateCrossConnectService(client *client.Client, ctx context.Context) PccsService {
	return &pccsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *pccsService) List(params ListQueryParams) (PrivateCrossConnects, *Response, error) {
	req := s.client.PrivateCrossConnectsApi.PccsGet(s.context)
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
	dcs, res, err := s.client.PrivateCrossConnectsApi.PccsGetExecute(req)
	return PrivateCrossConnects{dcs}, &Response{*res}, err
}

func (s *pccsService) Get(pccId string, params QueryParams) (*PrivateCrossConnect, *Response, error) {
	req := s.client.PrivateCrossConnectsApi.PccsFindById(s.context, pccId)
	pcc, res, err := s.client.PrivateCrossConnectsApi.PccsFindByIdExecute(req)
	return &PrivateCrossConnect{pcc}, &Response{*res}, err
}

func (s *pccsService) GetPeers(pccId string) (*[]Peer, *Response, error) {
	peers := make([]Peer, 0)
	req := s.client.PrivateCrossConnectsApi.PccsFindById(s.context, pccId)
	pcc, res, err := s.client.PrivateCrossConnectsApi.PccsFindByIdExecute(req)
	if err != nil {
		return nil, nil, err
	}
	if properties, ok := pcc.GetPropertiesOk(); ok && properties != nil {
		if ps, ok := properties.GetPeersOk(); ok && ps != nil {
			for _, p := range *ps {
				peers = append(peers, Peer{p})
			}
		}
	}
	return &peers, &Response{*res}, err
}

func (s *pccsService) Create(u PrivateCrossConnect, params QueryParams) (*PrivateCrossConnect, *Response, error) {
	req := s.client.PrivateCrossConnectsApi.PccsPost(s.context).Pcc(u.PrivateCrossConnect)
	pcc, res, err := s.client.PrivateCrossConnectsApi.PccsPostExecute(req)
	return &PrivateCrossConnect{pcc}, &Response{*res}, err
}

func (s *pccsService) Update(pccId string, input PrivateCrossConnectProperties, params QueryParams) (*PrivateCrossConnect, *Response, error) {
	req := s.client.PrivateCrossConnectsApi.PccsPatch(s.context, pccId).Pcc(input.PrivateCrossConnectProperties)
	pcc, res, err := s.client.PrivateCrossConnectsApi.PccsPatchExecute(req)
	return &PrivateCrossConnect{pcc}, &Response{*res}, err
}

func (s *pccsService) Delete(pccId string, params QueryParams) (*Response, error) {
	req := s.client.PrivateCrossConnectsApi.PccsDelete(s.context, pccId)
	res, err := s.client.PrivateCrossConnectsApi.PccsDeleteExecute(req)
	return &Response{*res}, err
}
