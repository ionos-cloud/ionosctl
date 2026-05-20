package waitfor

import (
	"context"
	"errors"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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

// Deprecated: WatchStateProgress is a legacy watcher. New code should rely on
// globalwait.WaitAndRerender. Only remaining caller: WaitForState (promote-volume path).
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
