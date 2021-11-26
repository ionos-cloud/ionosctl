package resources

import (
	"errors"
	"fmt"

	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
)

type Client struct {
	sdkgoauth.APIClient
}

type ClientConfig struct {
	sdkgoauth.Configuration
}

// ClientService is a wrapper around sdkgoauth.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *sdkgoauth.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(name, pwd, token, hostUrl string) (ClientService, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}
	clientConfig := sdkgoauth.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("ionos-cloud-sdk-go-auth-v%v-cli", sdkgoauth.Version)
	return &clientService{
		client: sdkgoauth.NewAPIClient(clientConfig),
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
