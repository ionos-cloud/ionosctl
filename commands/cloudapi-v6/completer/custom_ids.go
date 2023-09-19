/*
This is used for supporting completion in the CLI.
This is used for custom resources - filtered based on location, type, etc.
*/
package completer

import (
	"context"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func ImagesIdsCustom(_ io.Writer, params resources.ListQueryParams) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	imageSvc := resources.NewImageService(client, context.TODO())
	images, _, err := imageSvc.List(params)
	if err != nil {
		die.Die(err.Error())
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

func ServersIdsCustom(_ io.Writer, datacenterId string, params resources.ListQueryParams) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	serverSvc := resources.NewServerService(client, context.TODO())
	servers, _, err := serverSvc.List(datacenterId, params)
	if err != nil {
		die.Die(err.Error())
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
