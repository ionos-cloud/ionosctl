package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/auth-v1/resources"
	"github.com/spf13/viper"
)

func TokensIds(outErr io.Writer) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	tokenSvc := resources.NewTokenService(client, context.TODO())
	tokens, _, err := tokenSvc.List(0)
	clierror.CheckError(err, outErr)
	tokenIds := make([]string, 0)
	if items, ok := tokens.Tokens.GetTokensOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				tokenIds = append(tokenIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return tokenIds
}

// Get Client for Completion Functions
func getClient() (*resources.Client, error) {
	if err := config.Load(); err != nil {
		return nil, err
	}
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}
