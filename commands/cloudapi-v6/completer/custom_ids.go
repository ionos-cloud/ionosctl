/*
This is used for supporting completion in the CLI.
This is used for custom resources - filtered based on location, type, etc.
*/
package completer

import (
	"context"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func ImagesIdsCustom(outErr io.Writer, params resources.ListQueryParams) []string {
	client, err := client2.Get()
	clierror.CheckErrorAndDie(err, outErr)
	imageSvc := resources.NewImageService(client, context.TODO())
	images, _, err := imageSvc.List(params)
	clierror.CheckErrorAndDie(err, outErr)
	imgsIds := make([]string, 0)
	if items, ok := images.Images.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				imgsIds = append(imgsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return imgsIds
}

func ServersIdsCustom(outErr io.Writer, datacenterId string, params resources.ListQueryParams) []string {
	client, err := client2.Get()
	clierror.CheckErrorAndDie(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	servers, _, err := serverSvc.List(datacenterId, params)
	clierror.CheckErrorAndDie(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := servers.Servers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
