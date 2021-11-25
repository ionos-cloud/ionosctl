package resources

import (
	"errors"
	"strings"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

const DefaultBasePath = "/cloudapi/databases/postgres"

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
	if hostUrl == "" {
		return nil, errors.New("host-url incorrect")
	}
	if !strings.HasSuffix(hostUrl, DefaultBasePath) {
		hostUrl += DefaultBasePath
	}
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := &sdkgo.Configuration{
		Username: name,
		Password: pwd,
		Token:    token,
		Servers: sdkgo.ServerConfigurations{
			sdkgo.ServerConfiguration{
				URL: hostUrl,
			},
		},
	}
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
