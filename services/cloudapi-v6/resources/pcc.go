package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
)

type PrivateCrossConnect struct {
	compute.PrivateCrossConnect
}

type PrivateCrossConnectProperties struct {
	compute.PrivateCrossConnectProperties
}

type PrivateCrossConnects struct {
	compute.PrivateCrossConnects
}

type Peer struct {
	compute.Peer
}

// PccsService is a wrapper around compute.PrivateCrossConnect
type PccsService interface {
	List(params ListQueryParams) (PrivateCrossConnects, *Response, error)
	Get(pccId string, params QueryParams) (*PrivateCrossConnect, *Response, error)
	GetPeers(pccId string) (*[]Peer, *Response, error)
	Create(u PrivateCrossConnect, params QueryParams) (*PrivateCrossConnect, *Response, error)
	Update(pccId string, input PrivateCrossConnectProperties, params QueryParams) (*PrivateCrossConnect, *Response, error)
	Delete(pccId string, params QueryParams) (*Response, error)
}

type pccsService struct {
	client  *compute.APIClient
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
