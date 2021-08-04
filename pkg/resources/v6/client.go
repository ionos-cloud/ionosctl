package v6

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	if hostUrl == "" {
		return nil, errors.New("host-url incorrect")
	}
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := &ionoscloud.Configuration{
		Username: name,
		Password: pwd,
		Token:    token,
		Servers: ionoscloud.ServerConfigurations{
			ionoscloud.ServerConfiguration{
				URL: fmt.Sprintf("%s%s", hostUrl, config.DefaultV6BasePath),
			},
		},
	}
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
