package resources

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/viper"
)

type Client struct {
	sdkgo.APIClient
}

type ClientConfig struct {
	sdkgo.Configuration
}

// ClientService is a wrapper around ionoscloud.APIClient
type ClientService interface {
	GetById() *Client
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
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", viper.GetString(config.CLIHttpUserAgent), clientConfig.UserAgent)
	return &clientService{
		client: sdkgo.NewAPIClient(clientConfig),
	}, nil
}

func (c clientService) GetById() *Client {
	return &Client{
		APIClient: *c.client,
	}
}

func (c clientService) GetConfig() *ClientConfig {
	return &ClientConfig{
		Configuration: *c.client.GetConfig(),
	}
}
