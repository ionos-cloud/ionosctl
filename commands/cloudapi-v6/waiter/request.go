package waiter

import (
	"errors"

	"github.com/ionos-cloud/ionosctl/pkg/core"
)

func RequestInterrogator(c *core.CommandConfig, requestId string) (status *string, message *string, err error) {
	reqStatus, _, err := c.CloudApiV6Services.Requests().GetStatus(requestId)
	if err != nil {
		return nil, nil, err
	}
	if reqStatus != nil {
		if metadata, ok := reqStatus.GetMetadataOk(); ok && metadata != nil {
			if s, ok := metadata.GetStatusOk(); ok && s != nil {
				status = s
			}
			if msg, ok := metadata.GetMessageOk(); ok && msg != nil {
				message = msg
			}
			return status, message, nil
		}
	}
	return nil, nil, errors.New("error getting request status")
}
