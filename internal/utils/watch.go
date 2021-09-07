package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"time"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

// State Status
const (
	stateDeployingStatus = "DEPLOYING"
	stateUpdatingStatus  = "UPDATING"
	stateActiveStatus    = "ACTIVE"
	stateAvailableStatus = "AVAILABLE"
	stateReadyStatus     = "READY"
	stateFailedStatus    = "FAILED"
)

// WatchStateProgress watches the state progress of a Resource until it completes with success: meaning ACTIVE or AVAILABLE or error.
func WatchStateProgress(ctx context.Context, c *core.CommandConfig, interrogator InterrogateStateFunc, resourceId string) (<-chan int, <-chan error) {
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
			if state == nil {
				errChan <- errors.New("error getting state")
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
			case stateActiveStatus, stateAvailableStatus, stateReadyStatus:
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
func WatchRequestProgress(ctx context.Context, c *core.CommandConfig, interrogator InterrogateRequestFunc, requestId string) (<-chan int, <-chan error) {
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

			status, message, err := interrogator(c, requestId)
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
				errChan <- errors.New(fmt.Sprintf("%s %s", *status, *message))
				return
			}
		}
	}()
	return progressChan, errChan
}
