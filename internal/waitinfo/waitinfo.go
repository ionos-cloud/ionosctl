package waitinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	mu       sync.RWMutex
	lastHref string
)

// Set accepts any struct, marshals it to JSON, and extracts the "href" field.
func Set(dataStruct interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	// Marshal the struct into JSON
	dataAsJson, err := json.Marshal(dataStruct)
	if err != nil {
		return fmt.Errorf("failed marshalling struct to JSON: %w", err)
	}

	// Define a map to unmarshal JSON and extract "href"
	var jsonData map[string]interface{}
	if err := json.Unmarshal(dataAsJson, &jsonData); err != nil {
		return fmt.Errorf("failed unmarshalling JSON: %w", err)
	}

	// Extract the "href" field if it exists
	if href, exists := jsonData["href"]; exists {
		if hrefStr, ok := href.(string); ok {
			lastHref = hrefStr
		}
	}

	return nil
}

// GetHref returns the last stored href value.
func GetHref() string {
	mu.RLock()
	defer mu.RUnlock()
	return lastHref
}

// FindAndExecuteGetCommand tries to find the equivalent 'get' command to the given commandParts
//
// e.g. equivalent get command for "ionosctl datacenter delete --datacenter-id ID" is "ionosctl datacenter get --datacenter-id ID"
//
// and then executes it, in the hopes that doing so we will set the waitinfo.lastHref field.
func FindAndExecuteGetCommand(root *cobra.Command, commandParts, setFlagsWithValues []string) error {
	if len(commandParts) < 2 {
		return fmt.Errorf("failed to retrieve equivalent get command for %s: "+
			"commandParts must have at least 2 elements", strings.Join(commandParts, " "))
	}

	getCommand := append(commandParts[1:len(commandParts)-1], "get")
	foundCmd, _, err := root.Find(getCommand)
	if err != nil {
		return fmt.Errorf("failed to retrieve equivalent get command for %s "+
			"and no pre-existing output to deduce href to wait on", strings.Join(commandParts, " "))
	}

	// Temporarily override os.Args for execution
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }() // Restore after execution

	// From 'setFlagsWithValues' keep only the flags (with their values)
	// that are also defined on the 'get' command.
	var commonFlagsWithValues []string
	for _, flagKV := range setFlagsWithValues {
		// flagKV is in the form "--flag=value"
		if strings.HasPrefix(flagKV, "--") {
			trimmed := flagKV[2:]
			parts := strings.SplitN(trimmed, "=", 2)
			flagName := parts[0]
			// If the 'get' command defines this flag, keep it.
			fmt.Println("looking for flag: ", flagName)
			if foundCmd.Flags().Lookup(flagName) != nil {
				commonFlagsWithValues = append(commonFlagsWithValues, flagKV)
			}
		}
	}

	fmt.Println("commonFlagsWithValues: ", commonFlagsWithValues)

	// Build new args: binary name + getCommand slice + common flags.
	newArgs := append(getCommand, commonFlagsWithValues...)
	os.Args = append([]string{os.Args[0]}, newArgs...) // Preserve binary name
	foundCmd.SetArgs(commonFlagsWithValues)

	err = foundCmd.Execute()
	if err != nil {
		return fmt.Errorf("failed to execute get command '%s' equivalent to '%s' "+
			"and no pre-existing output to deduce href to wait on", foundCmd.CommandPath(), strings.Join(commandParts, " "))
	}

	if HrefIsEmpty() {
		return fmt.Errorf("no href could be deduced from the output of get command '%s' equivalent to '%s', "+
			"and no pre-existing output to deduce href to wait on", foundCmd.CommandPath(), strings.Join(commandParts, " "))
	}
	return nil
}

func HrefIsEmpty() bool {
	return lastHref == ""
}
