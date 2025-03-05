package lastresponse

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/itchyny/gojq"
)

var (
	mu             sync.RWMutex
	lastDataAsJson string
)

// SetE accepts any value but expects it to be serializable to JSON.
// If the data is not serializable to JSON, it returns an error.
func SetE(data interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	j, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		return fmt.Errorf("failed marshalling last API response: %w", errMarshal)
	}
	lastDataAsJson = string(j)
	return nil
}

// Set stores the JSON string directly.
func Set(dataAsJson string) {
	mu.Lock()
	defer mu.Unlock()
	lastDataAsJson = dataAsJson
}

// Get returns the last stored source data as a JSON string.
func Get() string {
	mu.RLock()
	defer mu.RUnlock()
	return lastDataAsJson
}

// GetJQ allows querying the last stored source data with a jq query.
// It returns the result of the query as a JSON string.
// If the query fails, it returns an error.
func GetJQ(query string) (string, error) {
	mu.RLock()
	data := lastDataAsJson
	mu.RUnlock()

	q, err := gojq.Parse(query)
	if err != nil {
		return "", fmt.Errorf("failed parsing query: %w", err)
	}

	var input interface{}
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("failed unmarshalling last data: %w", err)
	}

	iter := q.Run(input)
	var results []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", fmt.Errorf("query error: %w", err)
		}
		results = append(results, v)
	}

	// If there is a single result, return it directly
	if len(results) == 1 {
		switch v := results[0].(type) {
		case string:
			return v, nil
		case float64, int, bool, nil:
			return fmt.Sprintf("%v", v), nil
		default:
			res, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(res), nil
		}
	}

	// Otherwise, return a JSON array
	res, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
