package v5

import (
	"context"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	List() (PrivateCrossConnects, *Response, error)
	Get(pccId string) (*PrivateCrossConnect, *Response, error)
	GetPeers(pccId string) (*[]Peer, *Response, error)
	Create(u PrivateCrossConnect) (*PrivateCrossConnect, *Response, error)
	Update(pccId string, input PrivateCrossConnectProperties) (*PrivateCrossConnect, *Response, error)
	Delete(pccId string) (*Response, error)
}

type pccsService struct {
	client  *Client
	context context.Context
}

var _ PccsService = &pccsService{}

func NewPrivateCrossConnectService(client *Client, ctx context.Context) PccsService {
	return &pccsService{
		client:  client,
		context: ctx,
	}
}

func (s *pccsService) List() (PrivateCrossConnects, *Response, error) {
	req := s.client.PrivateCrossConnectApi.PccsGet(s.context)
	dcs, res, err := s.client.PrivateCrossConnectApi.PccsGetExecute(req)
	return PrivateCrossConnects{dcs}, &Response{*res}, err
}

func (s *pccsService) Get(pccId string) (*PrivateCrossConnect, *Response, error) {
	req := s.client.PrivateCrossConnectApi.PccsFindById(s.context, pccId)
	pcc, res, err := s.client.PrivateCrossConnectApi.PccsFindByIdExecute(req)
	return &PrivateCrossConnect{pcc}, &Response{*res}, err
}

func (s *pccsService) GetPeers(pccId string) (*[]Peer, *Response, error) {
	peers := make([]Peer, 0)
	req := s.client.PrivateCrossConnectApi.PccsFindById(s.context, pccId)
	pcc, res, err := s.client.PrivateCrossConnectApi.PccsFindByIdExecute(req)
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

func (s *pccsService) Create(u PrivateCrossConnect) (*PrivateCrossConnect, *Response, error) {
	req := s.client.PrivateCrossConnectApi.PccsPost(s.context).Pcc(u.PrivateCrossConnect)
	pcc, res, err := s.client.PrivateCrossConnectApi.PccsPostExecute(req)
	return &PrivateCrossConnect{pcc}, &Response{*res}, err
}

func (s *pccsService) Update(pccId string, input PrivateCrossConnectProperties) (*PrivateCrossConnect, *Response, error) {
	req := s.client.PrivateCrossConnectApi.PccsPatch(s.context, pccId).Pcc(input.PrivateCrossConnectProperties)
	pcc, res, err := s.client.PrivateCrossConnectApi.PccsPatchExecute(req)
	return &PrivateCrossConnect{pcc}, &Response{*res}, err
}

func (s *pccsService) Delete(pccId string) (*Response, error) {
	req := s.client.PrivateCrossConnectApi.PccsDelete(s.context, pccId)
	_, res, err := s.client.PrivateCrossConnectApi.PccsDeleteExecute(req)
	return &Response{*res}, err
}
