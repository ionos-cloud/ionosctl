package core

import (
	"encoding/json"
	"errors"
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

// ListAllLocations queries all locations concurrently when the API has more
// than one location and --location is not set, merging results into one view.
//
//   - text: merged table with "Location" as the first column.
//   - json: items from all locations merged under "items", each stamped with a
//     "location" field.
//   - api-json: array of per-location responses, each with a "location" field.
//
// For single-location APIs, non-regional commands, or when --location is set,
// it falls back to single-location behavior: the raw response is printed
// unchanged (no Location column, no merging, no array wrapping).
//
// The fetchFn receives a [shared.Configuration] for the target location URL.
// It must create its own SDK client from the config and execute the API call.
func (c *CommandConfig) ListAllLocations(
	columns []table.Column,
	fetchFn func(cfg *shared.Configuration) (any, error),
) error {
	locations, templateURL, productNames, found := findRegionalConfig(c.Command.Command)

	// Single-location behavior (raw passthrough, no Location column, no array
	// wrapping, no key stripping) when: no regional config, the API exposes a
	// single location (nothing to aggregate), or --location was set explicitly.
	// This keeps single-location APIs (DNS, CDN, Cert) byte-for-byte compatible
	// with pre-regional output.
	if !found || len(locations) <= 1 || c.Command.Command.Flags().Changed(constants.FlagLocation) {
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

// locationStampedItems extracts the "items" array from a collection response
// and stamps each object item with a top-level "location" field, so provenance
// survives the merge into a single {"items": [...]}. Non-object items are
// passed through unchanged. Returns an empty slice when the response has no
// "items" array.
func locationStampedItems(data any, location string) ([]any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	items, ok := m["items"].([]any)
	if !ok {
		return nil, nil
	}
	for _, item := range items {
		if obj, ok := item.(map[string]any); ok {
			obj["location"] = location
		}
	}
	return items, nil
}

// regionalLegacyJSON merges items from all locations into {"items": [...]}.
// Each item is stamped with its source location. This matches the
// single-location -o json shape (top-level "items"), with the added field.
func (c *CommandConfig) regionalLegacyJSON(results []locResult) error {
	if viper.GetBool(constants.ArgQuiet) {
		return nil
	}

	allItems := make([]any, 0)
	for _, r := range results {
		if r.err != nil {
			continue
		}
		items, err := locationStampedItems(r.data, r.location)
		if err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to process response from %s: %v\n", r.location, err)
			continue
		}
		allItems = append(allItems, items...)
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

// annotateLocation returns data with a top-level "location" field added, so
// per-location API responses carry their provenance in aggregated api-json
// output. It round-trips through JSON to avoid mutating the typed SDK response.
// If the response is not a JSON object (unexpected for collections), it is
// returned unmodified rather than dropped.
func annotateLocation(data any, location string) (any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		var raw any
		if err := json.Unmarshal(b, &raw); err != nil {
			return nil, err
		}
		return raw, nil
	}
	m["location"] = location
	return m, nil
}

// regionalAPIJSON outputs an array of per-location API responses. Each element
// is the API response for that location, annotated with a "location" field so
// provenance survives aggregation. Empty locations are kept (they reflect a
// location that was queried and returned no items).
func (c *CommandConfig) regionalAPIJSON(results []locResult) error {
	if viper.GetBool(constants.ArgQuiet) {
		return nil
	}

	rawResponses := make([]any, 0, len(results))
	for _, r := range results {
		if r.err != nil {
			continue
		}
		annotated, err := annotateLocation(r.data, r.location)
		if err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "WARN: failed to process response from %s: %v\n", r.location, err)
			continue
		}
		rawResponses = append(rawResponses, annotated)
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

// RunForAllLocations invokes fn once per allowed location (sequentially) for
// multi-location regional APIs when --location is not set, so bulk operations
// such as `delete --all` span every location the same way `list` does. Errors
// from individual locations are aggregated; a failure in one location does not
// stop the others.
//
// For non-regional commands, single-location APIs, or when --location is set
// explicitly, fn runs exactly once against the resolved single-location config.
//
// fn receives a per-location [shared.Configuration] and the location label; it
// must build its own SDK client from that config (do not use the global
// client.Must() singleton, which is bound to a single location).
func (c *CommandConfig) RunForAllLocations(fn func(cfg *shared.Configuration, location string) error) error {
	locations, templateURL, productNames, found := findRegionalConfig(c.Command.Command)

	if !found || len(locations) <= 1 || c.Command.Command.Flags().Changed(constants.FlagLocation) {
		loc := ""
		switch {
		case c.Command.Command.Flags().Changed(constants.FlagLocation):
			loc, _ = c.Command.Command.Flags().GetString(constants.FlagLocation)
		case len(locations) == 1:
			loc = locations[0]
		}
		return fn(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)), loc)
	}

	// Make the wider blast radius visible: a bulk op with no --location now
	// spans every location. Without this line, an existing `--all --force`
	// script that used to touch only the default location would silently act
	// across all of them.
	if !viper.GetBool(constants.ArgQuiet) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(),
			"Operating across all %d locations: %s (use --location to target a single one)\n",
			len(locations), strings.Join(locations, ", "))
	}

	var errs []error
	for _, loc := range locations {
		url := findOverridenURL(c.Command.Command, productNames, templateURL, loc)
		if err := fn(client.NewRegionalConfig(url), loc); err != nil {
			errs = append(errs, fmt.Errorf("location %s: %w", loc, err))
		}
	}
	return errors.Join(errs...)
}

// requireExplicitLocation returns an error when the command targets a regional
// API with more than one allowed location and --location was not explicitly
// provided. Resource IDs (and target locations for writes) are location-scoped
// and cannot be inferred, so single-location operations must know which
// location to target. Otherwise they silently default to the first allowed
// location, producing confusing 404s when the resource lives elsewhere.
//
// When the API exposes only a single location there is no ambiguity, so the
// default applies and no error is returned.
//
// Returns nil when the command is not regional, has a single location, or when
// --location was set.
func requireExplicitLocation(cmd *cobra.Command) error {
	locations, _, _, found := findRegionalConfig(cmd)
	if !found || len(locations) <= 1 || cmd.Flags().Changed(constants.FlagLocation) {
		return nil
	}
	return fmt.Errorf(
		"--location is required here: resources are location-scoped and cannot be inferred. "+
			"Set --location to one of: %s",
		strings.Join(locations, ", "),
	)
}

// RequireExplicitLocation enforces that --location is set for location-scoped
// (single-resource) operations on regional APIs. Call it from PreCmdRun of
// commands that operate on a specific resource ID. It is a no-op for
// non-regional commands or when --location is already set.
//
// List commands that fan out over all locations must NOT call this; they use
// [CommandConfig.ListAllLocations] instead.
func (c *PreCommandConfig) RequireExplicitLocation() error {
	return requireExplicitLocation(c.Command.Command)
}

// RequireExplicitLocation enforces that --location is set for location-scoped
// (single-resource) operations on regional APIs. This CommandConfig variant is
// for hybrid list commands that only take a single-location branch at runtime
// (e.g. when a specific parent ID is provided); call it inside that branch. It
// is a no-op for non-regional commands, single-location APIs, or when
// --location is already set.
func (c *CommandConfig) RequireExplicitLocation() error {
	return requireExplicitLocation(c.Command.Command)
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
