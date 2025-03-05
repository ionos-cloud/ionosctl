package lastresponse

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	mu       sync.RWMutex
	lastData map[string]interface{}
)

// SetE accepts any value but expects it to be a map[string]interface{}.
// It returns an error if the type is incorrect.
func SetE(data interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	j, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		return fmt.Errorf("failed marshalling last API response: %w", errMarshal)
	}
	errUnmarshal := json.Unmarshal(j, &lastData)
	if errUnmarshal != nil {
		return fmt.Errorf("failed unmarshalling last API response: %w", errUnmarshal)
	}
	return nil
}

func Set(data map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	lastData = data
}

// Get returns the last stored source data.
func Get() map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return lastData
}
