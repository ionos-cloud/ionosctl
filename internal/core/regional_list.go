package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	// AnnotationLocations is set by WithRegionalConfigOverride on the service
	// root command. It contains comma-separated allowed locations (e.g. "de/fra,de/txl").
	AnnotationLocations = "regional.locations"

	// AnnotationTemplateURL is the URL template with %s placeholder for location.
	AnnotationTemplateURL = "regional.templateURL"
)

type locResult struct {
	location string
	data     any
	err      error
}

// ListAllLocations queries all locations concurrently when --location is not
// explicitly set, merging results into a single table with a Location column.
//
// For text output: merged table with "Location" as the first column.
// For JSON/api-json output: array of untouched per-location API responses.
// When --location is explicitly set: single-location behavior (unchanged).
//
// The fetchFn receives a [shared.Configuration] for the target location URL.
// It must create its own SDK client from the config and execute the API call.
func (c *CommandConfig) ListAllLocations(
	columns []table.Column,
	fetchFn func(cfg *shared.Configuration) (any, error),
) error {
	locations, templateURL, found := findRegionalConfig(c.Command.Command)

	// No regional config or --location explicitly set: single-location behavior
	if !found || c.Command.Command.Flags().Changed(constants.FlagLocation) {
		cfg := client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl))
		data, err := fetchFn(cfg)
		if err != nil {
			return err
		}
		return c.Printer(columns).Prefix("items").Print(data)
	}

	// Multi-location: query all concurrently
	results := make([]locResult, len(locations))
	var wg sync.WaitGroup

	for i, loc := range locations {
		wg.Add(1)
		go func(i int, loc string) {
			defer wg.Done()

			normalizedLoc := strings.ReplaceAll(loc, "/", "-")
			url := fmt.Sprintf(templateURL, normalizedLoc)
			cfg := client.NewRegionalConfig(url)

			data, err := fetchFn(cfg)
			results[i] = locResult{location: loc, data: data, err: err}
		}(i, loc)
	}
	wg.Wait()

	// Check if all failed
	var lastErr error
	anySuccess := false
	for _, r := range results {
		if r.err != nil {
			lastErr = r.err
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to list from %s: %v\n", r.location, r.err)
		} else {
			anySuccess = true
		}
	}
	if !anySuccess {
		return fmt.Errorf("failed to list from any location: %w", lastErr)
	}

	// Determine output format
	format := viper.GetString(constants.ArgOutput)
	if format == "json" || format == "api-json" {
		return c.regionalJSON(results)
	}

	return c.regionalText(results, columns)
}

// regionalJSON outputs an array of per-location API responses.
// Both -o json and -o api-json produce the same format: an array of raw responses.
func (c *CommandConfig) regionalJSON(results []locResult) error {
	if viper.GetBool(constants.ArgQuiet) {
		return nil
	}

	var rawResponses []any
	for _, r := range results {
		if r.err == nil {
			rawResponses = append(rawResponses, r.data)
		}
	}

	var data any = rawResponses
	data, err := table.ApplyQueryFilter(data)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return c.Out(string(out)+"\n", nil)
}

// regionalText builds a merged table with Location as the first column.
func (c *CommandConfig) regionalText(results []locResult, columns []table.Column) error {
	// Prepend Location column
	allCols := append([]table.Column{{Name: "Location", Default: true}}, columns...)
	t := table.New(allCols)

	for _, r := range results {
		if r.err != nil {
			continue
		}

		// Extract rows from this location's response
		locTable := table.New(columns, table.WithPrefix("items"))
		if err := locTable.Extract(r.data); err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to parse response from %s: %v\n", r.location, err)
			continue
		}

		for _, row := range locTable.Rows() {
			row["Location"] = r.location
			t.AppendRow(row)
		}
	}

	return c.Out(t.Render(table.ResolveCols(allCols, c.Cols())))
}

// findRegionalConfig walks parent commands to find regional annotations
// set by [WithRegionalConfigOverride].
func findRegionalConfig(cmd *cobra.Command) (locations []string, templateURL string, found bool) {
	for c := cmd; c != nil; c = c.Parent() {
		locs, hasLocs := c.Annotations[AnnotationLocations]
		tmpl, hasTmpl := c.Annotations[AnnotationTemplateURL]
		if hasLocs && hasTmpl {
			return strings.Split(locs, ","), tmpl, true
		}
	}
	return nil, "", false
}
