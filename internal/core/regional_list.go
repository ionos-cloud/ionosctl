package core

import (
	"encoding/json"
	"fmt"
	"sort"
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

	// AnnotationProductNames is the comma-separated product names used to look
	// up per-location config-file URL overrides (mirrors WithRegionalConfigOverride).
	AnnotationProductNames = "regional.productNames"
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
	locations, templateURL, productNames, found := findRegionalConfig(c.Command.Command)

	// No regional config or --location explicitly set: single-location behavior
	if !found || c.Command.Command.Flags().Changed(constants.FlagLocation) {
		cfg := client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl))
		data, err := fetchFn(cfg)
		if err != nil {
			return err
		}
		return c.Printer(columns).Prefix("items").Print(data)
	}

	// Build all configs before spawning goroutines to avoid racing
	// on client.Must() singleton (getters.go URL-change rebuild is not synchronized).
	type locConfig struct {
		location string
		cfg      *shared.Configuration
	}
	configs := make([]locConfig, len(locations))
	for i, loc := range locations {
		// Resolve per-location URL honoring overrides (--api-url, env var,
		// per-location config-file override), falling back to the template.
		url := findOverridenURL(c.Command.Command, productNames, templateURL, loc)
		configs[i] = locConfig{location: loc, cfg: client.NewRegionalConfig(url)}
	}

	// Query all locations concurrently
	results := make([]locResult, len(locations))
	var wg sync.WaitGroup

	for i, lc := range configs {
		wg.Add(1)
		go func(i int, lc locConfig) {
			defer wg.Done()
			data, err := fetchFn(lc.cfg)
			results[i] = locResult{location: lc.location, data: data, err: err}
		}(i, lc)
	}
	wg.Wait()

	// Collect errors, deduplicate identical messages
	var lastErr error
	anySuccess := false
	errCounts := map[string][]string{} // error message to list of locations
	for _, r := range results {
		if r.err != nil {
			lastErr = r.err
			msg := r.err.Error()
			errCounts[msg] = append(errCounts[msg], r.location)
		} else {
			anySuccess = true
		}
	}

	if !anySuccess {
		// All failed : return single error, no warnings
		if len(errCounts) == 1 {
			return fmt.Errorf("failed to list from all locations: %w", lastErr)
		}
		// Multiple distinct errors : join them (sorted for stable output)
		var parts []string
		for msg, locs := range errCounts {
			sort.Strings(locs)
			parts = append(parts, fmt.Sprintf("%s: %s", strings.Join(locs, ", "), msg))
		}
		sort.Strings(parts)
		return fmt.Errorf("failed to list from all locations:\n  %s", strings.Join(parts, "\n  "))
	}

	// Partial failure : warn only for failed locations (sorted for stable output)
	if len(errCounts) > 0 && !viper.GetBool(constants.ArgQuiet) {
		stderr := c.Command.Command.ErrOrStderr()
		var warns []string
		for msg, locs := range errCounts {
			sort.Strings(locs)
			warns = append(warns, fmt.Sprintf("WARN: failed to list from %s: %s", strings.Join(locs, ", "), msg))
		}
		sort.Strings(warns)
		for _, w := range warns {
			fmt.Fprintln(stderr, w)
		}
	}

	// Determine output format
	format := viper.GetString(constants.ArgOutput)
	switch format {
	case "json":
		return c.regionalLegacyJSON(results)
	case "api-json":
		return c.regionalAPIJSON(results)
	}

	return c.regionalText(results, columns)
}

// regionalLegacyJSON merges items from all locations into {"items": [...]}.
// This matches the single-location -o json format (non-breaking).
func (c *CommandConfig) regionalLegacyJSON(results []locResult) error {
	if viper.GetBool(constants.ArgQuiet) {
		return nil
	}

	allItems := make([]any, 0)
	for _, r := range results {
		if r.err != nil {
			continue
		}
		b, err := json.Marshal(r.data)
		if err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to marshal response from %s: %v\n", r.location, err)
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to parse response from %s: %v\n", r.location, err)
			continue
		}
		if items, ok := m["items"].([]any); ok {
			allItems = append(allItems, items...)
		}
	}

	var data any = map[string]any{"items": allItems}
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

// regionalAPIJSON outputs an array of per-location API responses.
// Each element is the raw, untouched API response for that location.
func (c *CommandConfig) regionalAPIJSON(results []locResult) error {
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
			if !viper.GetBool(constants.ArgQuiet) {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to parse response from %s: %v\n", r.location, err)
			}
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
func findRegionalConfig(cmd *cobra.Command) (locations []string, templateURL string, productNames []string, found bool) {
	for c := cmd; c != nil; c = c.Parent() {
		locs, hasLocs := c.Annotations[AnnotationLocations]
		tmpl, hasTmpl := c.Annotations[AnnotationTemplateURL]
		if hasLocs && hasTmpl {
			var prods []string
			if p, ok := c.Annotations[AnnotationProductNames]; ok && p != "" {
				prods = strings.Split(p, ",")
			}
			return strings.Split(locs, ","), tmpl, prods, true
		}
	}
	return nil, "", nil, false
}
