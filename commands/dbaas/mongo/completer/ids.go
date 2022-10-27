package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/dbaas-mongo/resources"
	"github.com/spf13/viper"
	"io"
)

func MongoTemplateIds(outErr io.Writer) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	svc := resources.NewTemplatesService(client, context.TODO())
	templates, _, err := svc.List()
	clierror.CheckError(err, outErr)
	ids := make([]string, 0)
	if items, ok := templates.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ids = append(ids, *itemId)
			}
		}
	} else {
		return nil
	}
	return ids
}

// Get Client for Completion Functions
// TODO: we should use a "client" package... this makes no sense to be in EVERY SINGLE completions package
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
