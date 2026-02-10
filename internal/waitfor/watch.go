package waitfor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// State Status
const (
	stateDeployingStatus  = "DEPLOYING"
	stateUpdatingStatus   = "UPDATING"
	stateBusyStatus       = "BUSY"
	stateActiveStatus     = "ACTIVE"
	stateAvailableStatus  = "AVAILABLE"
	stateReadyStatus      = "READY"
	stateDoneStatus       = "DONE"
	stateFailedStatus     = "FAILED"
	stateDestroyingStatus = "DESTROYING"
)

// WatchStateProgress watches the state progress of a Resource until it completes with success: meaning ACTIVE or AVAILABLE or error.
func WatchStateProgress(ctx context.Context, c *core.CommandConfig, interrogator InterrogateStateFunc, resourceId string) (<-chan int, <-chan error) {
	errChan := make(chan error, 1)
	progressChan := make(chan int)
	go func() {
		defer close(errChan)
		defer close(progressChan)
		ticker := time.NewTicker(pollTime)
		defer ticker.Stop()
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
				// perhaps the resource is being created, let it wait longer
				sendingProgress(0)
			} else {
				switch *state {
				case stateDeployingStatus:
					sendingProgress(5)
					break
				case stateUpdatingStatus, stateBusyStatus:
					sendingProgress(50)
					break
				case stateActiveStatus, stateAvailableStatus, stateReadyStatus, stateDoneStatus:
					sendingProgress(100)
					errChan <- nil
					return
				case stateFailedStatus, stateDestroyingStatus:
					errChan <- errors.New(failed)
					return
				}
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
		defer ticker.Stop()
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

// WatchDeletionProgress watches the deletion progress of a Resource until it completes with success: returning 404 http response status code.
func WatchDeletionProgress(ctx context.Context, c *core.CommandConfig, interrogator InterrogateDeletionFunc, resourceId string) (<-chan int, <-chan error) {
	errChan := make(chan error, 1)
	progressChan := make(chan int)
	go func() {
		defer close(errChan)
		defer close(progressChan)
		ticker := time.NewTicker(pollTime)
		defer ticker.Stop()
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

			httpResponseCode, err := interrogator(c, resourceId)
			if err != nil {
				errChan <- err
				return
			}
			if httpResponseCode == nil {
				errChan <- errors.New("error getting http response status code")
				return
			}

			// Check Resource Http Response Status Code
			// Send Progress, Send Error if any
			switch *httpResponseCode {
			case 200:
				sendingProgress(1)
				break
			case 202:
				sendingProgress(50)
				break
			case 404:
				sendingProgress(100)
				errChan <- nil
				return
			case 400, 401, 403, 500:
				errChan <- errors.New(failed)
				return
			}
		}
	}()
	return progressChan, errChan
}
