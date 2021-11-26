package resources

import (
	"errors"
	"fmt"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Client struct {
	ionoscloud.APIClient
}

type ClientConfig struct {
	ionoscloud.Configuration
}

// ClientService is a wrapper around ionoscloud.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *ionoscloud.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(name, pwd, token, hostUrl string) (ClientService, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := ionoscloud.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("ionos-cloud-sdk-go-v%v-cli", ionoscloud.Version)
	return &clientService{
		client: ionoscloud.NewAPIClient(clientConfig),
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
