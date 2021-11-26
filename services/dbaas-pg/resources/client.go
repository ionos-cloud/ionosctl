package resources

import (
	"errors"
	"fmt"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type Client struct {
	sdkgo.APIClient
}

type ClientConfig struct {
	sdkgo.Configuration
}

// ClientService is a wrapper around ionoscloud.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *sdkgo.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(name, pwd, token, hostUrl string) (ClientService, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := sdkgo.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("ionos-cloud-sdk-go-dbaas-postgres-%v-cli", sdkgo.Version)
	return &clientService{
		client: sdkgo.NewAPIClient(clientConfig),
	}, nil
}

func (c clientService) Get() *Client {
	return &Client{
		APIClient: *c.client,
	}
}

func (c clientService) GetConfig() *ClientConfig {
	return &ClientConfig{
		Configuration: *c.client.GetConfig(),
	}
}
