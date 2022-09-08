package resources

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	dp "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/viper"
)

type Client struct {
	dp.APIClient
}

type ClientConfig struct {
	dp.Configuration
}

// ClientService is a wrapper around ionoscloud.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *dp.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(name, pwd, token, hostUrl string) (ClientService, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := dp.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", viper.GetString(config.CLIHttpUserAgent), clientConfig.UserAgent)
	return &clientService{
		client: dp.NewAPIClient(clientConfig),
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
