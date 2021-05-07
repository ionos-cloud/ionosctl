package utils

import (
	"context"
	"errors"
	"time"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

// State Status
const (
	stateDeployingStatus = "DEPLOYING"
	stateUpdatingStatus  = "UPDATING"
	stateActiveStatus    = "ACTIVE"
	stateAvailableStatus = "AVAILABLE"
	stateFailedStatus    = "FAILED"
)

// WatchStateProgress watches the state progress of a Resource until it completes with success: meaning ACTIVE or AVAILABLE or error.
func WatchStateProgress(ctx context.Context, c *builder.CommandConfig, interrogator InterrogateStateFunc, resourceId string) (<-chan int, <-chan error) {
	errChan := make(chan error, 1)
	progressChan := make(chan int)
	go func() {
		defer close(errChan)
		defer close(progressChan)
		ticker := time.NewTicker(pollTime)
		sendingProgress := func(p int) {
			select {
			case progressChan <- p:
				break
			default:
				break
			}
		}
		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case <-ticker.C:
				break
			}

			state, err := interrogator(c, resourceId)
			if err != nil {
				errChan <- err
				return
			}

			// Check Resource State
			// Send Progress, Send Error if any
			switch *state {
			case stateDeployingStatus:
				sendingProgress(1)
				break
			case stateUpdatingStatus:
				sendingProgress(50)
				break
			case stateActiveStatus:
				sendingProgress(100)
				errChan <- nil
				return
			case stateAvailableStatus:
				sendingProgress(100)
				errChan <- nil
				return
			case stateFailedStatus:
				errChan <- errors.New(failed)
				return
			}
		}
	}()
	return progressChan, errChan
}

// WatchRequestProgress watches the status progress of a Request until it completes with success: meaning DONE or error.
func WatchRequestProgress(ctx context.Context, c *builder.CommandConfig, requestId string) (<-chan int, <-chan error) {
	errChan := make(chan error, 1)
	progressChan := make(chan int)
	go func() {
		defer close(errChan)
		defer close(progressChan)
		ticker := time.NewTicker(pollTime)
		sendingProgress := func(p int) {
			select {
			case progressChan <- p:
				break
			default:
				break
			}
		}
		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case <-ticker.C:
				break
			}

			reqStatus, _, err := c.Requests().GetStatus(requestId)
			if err != nil {
				errChan <- err
				return
			}
			status, err := getRequestStatus(reqStatus)
			if err != nil {
				errChan <- err
				return
			}

			// Check Resource State
			// Send Progress, Send Error if any
			switch *status {
			case ionoscloud.RequestStatusQueued:
				sendingProgress(1)
				break
			case ionoscloud.RequestStatusRunning:
				sendingProgress(50)
				break
			case ionoscloud.RequestStatusDone:
				sendingProgress(100)
				errChan <- nil
				return
			case ionoscloud.RequestStatusFailed:
				errChan <- errors.New(failed)
				return
			}
		}
	}()
	return progressChan, errChan
}

func getRequestStatus(reqStatus *resources.RequestStatus) (*string, error) {
	if metadata, ok := reqStatus.GetMetadataOk(); ok && metadata != nil {
		if status, ok := metadata.GetStatusOk(); ok && status != nil {
			return status, nil
		}
	}
	return nil, errors.New("error getting request status")
}
