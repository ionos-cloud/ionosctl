package resources

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	cert "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/viper"
)

type Client struct {
	cert.APIClient
}

type ClientConfig struct {
	cert.Configuration
}

// ClientService is a wrapper around ionoscloud.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *cert.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(name, pwd, token, hostUrl string) (ClientService, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := cert.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), clientConfig.UserAgent)
	return &clientService{
		client: cert.NewAPIClient(clientConfig),
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
