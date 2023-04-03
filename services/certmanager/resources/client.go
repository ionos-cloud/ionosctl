package resources

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/cert"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

type Client struct {
	sdkgo.APIClient
}

type ClientConfig struct {
	shared.Configuration
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
	clientConfig := shared.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), clientConfig.UserAgent)
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
