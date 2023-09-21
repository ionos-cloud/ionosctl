/*
This is used for supporting completion in the CLI.
This is used for custom resources - filtered based on location, type, etc.
*/
package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func ImagesIdsCustom(outErr io.Writer, params resources.ListQueryParams) []string {
	imageSvc := resources.NewImageService(client.Must(), context.Background())
	images, _, err := imageSvc.List(params)
	if err != nil {
		return nil
	}
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
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	servers, _, err := serverSvc.List(datacenterId, params)
	if err != nil {
		return nil
	}
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
