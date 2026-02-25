// Package table provides a unified Table abstraction for formatting CLI output.
//
// It replaces the scattered jsontabwriter + tabheaders + jsonpaths + resource2table
// packages with a single, cohesive API. Column definitions carry their own JSON extraction
// paths, default visibility, and optional format functions — eliminating the need for
// separate "preconverted" output paths.
//
// Basic usage:
//
//	var datacenterCols = []table.Column{
//	    {Name: "DatacenterId", JSONPath: "id", Default: true},
//	    {Name: "Name", JSONPath: "properties.name", Default: true},
//	    {Name: "State", JSONPath: "metadata.state", Default: true},
//	    {Name: "Description", JSONPath: "properties.description"},
//	}
//
//	// In command execution:
//	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
//	return c.Out(table.Sprint(datacenterCols, dc, cols))
//
// For resources with computed columns (replaces resource2table converters):
//
//	var clusterCols = []table.Column{
//	    {Name: "ClusterId", JSONPath: "id", Default: true},
//	    {Name: "MaintenanceWindow", Default: true, Format: func(item map[string]any) any {
//	        mw := table.Navigate(item, "properties.maintenanceWindow")
//	        if mw == nil { return "" }
//	        m := mw.(map[string]any)
//	        return fmt.Sprintf("%s %s", m["dayOfTheWeek"], m["time"])
//	    }},
//	}
//
// For more control (editing cells, re-rendering with fresh data):
//
//	t := table.New(datacenterCols)
//	t.Extract(dc)
//	t.SetCell(0, "State", "AVAILABLE")
//	return c.Out(t.Render(table.ResolveCols(datacenterCols, userCols)))
package table

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/Jeffail/gabs/v2"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/viper"
)

// Column defines a single column in a table.
type Column struct {
	Name     string     // Display name / header (e.g., "ServerId", "Name", "State")
	JSONPath string     // Extraction path (e.g., "id", "properties.name", "metadata.state")
	Default  bool       // Show in default output (when user doesn't specify --cols)
	Format   FormatFunc // Optional: custom value transformer (overrides JSONPath extraction)
}

// FormatFunc transforms a value for text display.
// It receives the full raw JSON item as a map, allowing access to any field.
// This enables computed columns like combining maintenance window day + time.
// Return nil or "" to indicate no value (column will be empty for this row).
type FormatFunc func(item map[string]any) any

// BeforeRenderFunc is a hook called before rendering. If it returns false,
// rendering is suppressed (returns "", nil). This is used by the --wait flag
// to defer output until the resource reaches a terminal state.
type BeforeRenderFunc func(t *Table, visibleCols []string) bool

// BeforeRender is a package-level hook that, if set, is called before each
// Render invocation. Set this in init() or early startup to integrate with
// features like global --wait.
var BeforeRender BeforeRenderFunc

// Table holds column definitions, extracted row data, and original source data.
// Create with New, populate with Extract, format with Render.
type Table struct {
	columns []Column
	prefix  string
	rows    []map[string]any
	raw     any
}

// New creates a Table with the given column definitions and options.
func New(columns []Column, opts ...Option) *Table {
	t := &Table{columns: columns}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Option configures Table creation.
type Option func(*Table)

// WithPrefix sets the JSONPath prefix for navigating to the data array root.
// For example, "items" navigates into the "items" array of a list response.
func WithPrefix(prefix string) Option {
	return func(t *Table) { t.prefix = prefix }
}

// Extract populates the table rows from source data (SDK struct, slice, or map).
// It stores the raw data for JSON output and extracts formatted rows for text output.
// Can be called multiple times to replace data (e.g., after re-fetching for --wait).
func (t *Table) Extract(sourceData any) error {
	if sourceData == nil {
		return fmt.Errorf("source data cannot be nil")
	}

	t.raw = sourceData

	b, err := json.Marshal(sourceData)
	if err != nil {
		return fmt.Errorf("marshal source data: %w", err)
	}

	// Handle empty arrays
	if reflect.DeepEqual(b, []byte{'[', ']'}) {
		t.rows = nil
		return nil
	}

	parsed, err := gabs.ParseJSON(b)
	if err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	items, err := navigateToRoot(parsed, t.prefix, sourceData)
	if err != nil {
		return err
	}

	t.rows = make([]map[string]any, 0, len(items))
	for _, item := range items {
		rawItem := containerToMap(item)
		row := make(map[string]any)

		for _, col := range t.columns {
			if col.Format != nil {
				row[col.Name] = col.Format(rawItem)
			} else if col.JSONPath != "" {
				row[col.Name] = extractValue(item, col.Name, col.JSONPath)
			}
		}

		t.rows = append(t.rows, row)
	}

	return nil
}

// Rows returns the extracted row data. Each map represents a row with column names as keys.
func (t *Table) Rows() []map[string]any {
	return t.rows
}

// Raw returns the original source data.
func (t *Table) Raw() any {
	return t.raw
}

// SetCell sets a value for a specific column in a specific row.
// Useful for post-extraction adjustments.
func (t *Table) SetCell(row int, col string, value any) {
	if row >= 0 && row < len(t.rows) {
		t.rows[row][col] = value
	}
}

// Render produces the final formatted output string.
// It respects --output, --quiet, --query, and --no-headers flags.
// visibleCols specifies which columns to display; use ResolveCols to compute this
// from user-supplied --cols values.
//
// If the package-level BeforeRender hook is set and returns false, output is
// suppressed (returns "", nil).
func (t *Table) Render(visibleCols []string) (string, error) {
	if viper.GetBool(constants.ArgQuiet) {
		return "", nil
	}

	if BeforeRender != nil && !BeforeRender(t, visibleCols) {
		return "", nil
	}

	format := viper.GetString(constants.ArgOutput)
	switch format {
	case "text":
		if viper.GetString(constants.FlagQuery) != "" {
			return "", fmt.Errorf("JMESPath filtering (--query) is not supported with text output. Use -o api-json or json format instead")
		}
		return t.renderText(visibleCols), nil
	case "json":
		return t.renderLegacyJSON()
	case "api-json":
		return t.renderAPIJSON()
	default:
		return "", fmt.Errorf("invalid format: %s", format)
	}
}

// Sprint is the single convenience entry point: Extract + Render → string.
//
//	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
//	return c.Out(table.Sprint(serverCols, data, cols, table.WithPrefix("items")))
func Sprint(columns []Column, sourceData any, userCols []string, opts ...Option) (string, error) {
	t := New(columns, opts...)
	if err := t.Extract(sourceData); err != nil {
		return "", err
	}
	return t.Render(ResolveCols(columns, userCols))
}

// --- Column helpers (replaces tabheaders package) ---

// AllCols returns the names of all defined columns.
func AllCols(columns []Column) []string {
	names := make([]string, len(columns))
	for i, c := range columns {
		names[i] = c.Name
	}
	return names
}

// DefaultCols returns the names of columns marked as Default.
func DefaultCols(columns []Column) []string {
	var names []string
	for _, c := range columns {
		if c.Default {
			names = append(names, c.Name)
		}
	}
	return names
}

// ResolveCols resolves user-specified column names against defined columns.
// Handles the "all" keyword, case-insensitive matching, and falls back to defaults.
func ResolveCols(columns []Column, userCols []string) []string {
	all := AllCols(columns)
	defaults := DefaultCols(columns)

	if len(userCols) > 0 && userCols[0] == "all" {
		return all
	}

	if userCols == nil {
		return defaults
	}

	allLower := make([]string, len(all))
	for i, c := range all {
		allLower[i] = strings.ToLower(c)
	}

	var valid []string
	for _, uc := range userCols {
		for i, lc := range allLower {
			if strings.ToLower(uc) == lc {
				valid = append(valid, all[i])
				break
			}
		}
	}

	if len(valid) == 0 {
		return defaults
	}
	return valid
}

// ColsMessage generates help text for the --cols flag.
func ColsMessage(columns []Column) string {
	return fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", AllCols(columns))
}

// --- Utility ---

// Navigate extracts a nested value from a map using a dot-separated path.
// Useful inside FormatFunc implementations.
//
//	mw := table.Navigate(item, "properties.maintenanceWindow")
func Navigate(m map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var current any = m
	for _, part := range parts {
		cm, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = cm[part]
	}
	return current
}

// --- Internal rendering ---

func (t *Table) renderText(visibleCols []string) string {
	activeCols := eliminateEmptyCols(visibleCols, t.rows)
	if len(activeCols) == 0 {
		return ""
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 5, 0, 3, ' ', tabwriter.StripEscape)

	if !viper.GetBool(constants.ArgNoHeaders) {
		fmt.Fprintln(w, strings.Join(activeCols, "\t"))
	}

	for _, row := range t.rows {
		vals := make([]string, len(activeCols))
		for i, col := range activeCols {
			vals[i] = formatCellValue(row[col])
		}
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	}

	w.Flush()
	return buf.String()
}

func (t *Table) renderAPIJSON() (string, error) {
	data := t.raw
	data, err := applyQueryFilter(data)
	if err != nil {
		return "", err
	}
	return marshalJSON(data)
}

func (t *Table) renderLegacyJSON() (string, error) {
	data, err := legacyJSONWrap(t.raw)
	if err != nil {
		return "", err
	}
	data, err = applyQueryFilter(data)
	if err != nil {
		return "", err
	}
	return marshalJSON(data)
}

// --- Internal helpers ---

func navigateToRoot(parsed *gabs.Container, prefix string, sourceData any) ([]*gabs.Container, error) {
	if prefix == "" {
		if reflect.TypeOf(sourceData).Kind() == reflect.Slice {
			return parsed.Children(), nil
		}
		return []*gabs.Container{parsed}, nil
	}

	if !parsed.ExistsP(prefix) {
		return nil, nil
	}

	root := parsed.Path(prefix)
	children := root.Children()
	if children == nil {
		return nil, fmt.Errorf("prefix %q does not lead to an array", prefix)
	}

	// Handle nested arrays (flatten)
	var flattened []*gabs.Container
	for _, child := range children {
		if _, ok := child.Data().([]any); ok {
			flattened = append(flattened, child.Children()...)
		} else {
			return children, nil
		}
	}
	return flattened, nil
}

func containerToMap(c *gabs.Container) map[string]any {
	if m, ok := c.Data().(map[string]any); ok {
		return m
	}
	// Fallback for non-map data
	b, _ := json.Marshal(c.Data())
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}

func extractValue(item *gabs.Container, colName, path string) any {
	data := item.Path(path)
	// "href as parent column" convention: when a column uses the href path
	// but isn't named "href", extract the parent resource ID from the URL.
	if strings.Contains(path, "href") && strings.ToLower(colName) != "href" {
		return parentIDFromHref(data.Data())
	}
	return data.Data()
}

func parentIDFromHref(href any) string {
	s, ok := href.(string)
	if !ok {
		return ""
	}
	tokens := strings.Split(s, "/")
	if len(tokens) < 3 {
		return ""
	}
	// Parent ID is 2 tokens before the child ID in IONOS API hrefs
	return tokens[len(tokens)-3]
}

func eliminateEmptyCols(cols []string, rows []map[string]any) []string {
	var active []string
	for _, c := range cols {
		for _, row := range rows {
			v, ok := row[c]
			if !ok || v == nil {
				continue
			}
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Slice && rv.Len() == 0 {
				continue
			}
			if s, ok := v.(string); ok && s == "" {
				continue
			}
			active = append(active, c)
			break
		}
	}
	return active
}

func formatCellValue(v any) string {
	if v == nil {
		return ""
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice:
		parts := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			parts[i] = fmt.Sprintf("%v", rv.Index(i).Interface())
		}
		return strings.Join(parts, ", ")
	case reflect.Float64:
		f := v.(float64)
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		return fmt.Sprintf("%v", f)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func marshalJSON(data any) (string, error) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out) + "\n", nil
}

func applyQueryFilter(data any) (any, error) {
	if !viper.IsSet(constants.FlagQuery) {
		return data, nil
	}
	expr := viper.GetString(constants.FlagQuery)
	if expr == "" {
		return data, nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal for filter: %w", err)
	}
	var generic any
	if err := json.Unmarshal(b, &generic); err != nil {
		return nil, fmt.Errorf("unmarshal for filter: %w", err)
	}
	result, err := jmespath.Search(expr, generic)
	if err != nil {
		return nil, fmt.Errorf("failed applying filter %q: %w", expr, err)
	}
	return result, nil
}

// legacyJSONWrap wraps source data into the legacy {"items": [...]} format for -o json.
func legacyJSONWrap(sourceData any) (any, error) {
	raw, err := json.Marshal(sourceData)
	if err != nil {
		return nil, fmt.Errorf("marshal for legacy JSON: %w", err)
	}

	var temp any
	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, fmt.Errorf("unmarshal for legacy JSON: %w", err)
	}

	slice, ok := temp.([]any)
	if !ok {
		return temp, nil
	}

	merged := make([]any, 0)
	foundItemsSlice := false
	for _, elem := range slice {
		m, isMap := elem.(map[string]any)
		if !isMap {
			continue
		}
		itemsVal, hasItems := m["items"]
		if !hasItems {
			continue
		}
		itemsSlice, isSlice := itemsVal.([]any)
		if !isSlice {
			continue
		}
		foundItemsSlice = true
		merged = append(merged, itemsSlice...)
	}

	if foundItemsSlice {
		return map[string]any{"items": merged}, nil
	}

	fallback := make([]any, 0, len(slice))
	for _, elem := range slice {
		m, isMap := elem.(map[string]any)
		if isMap {
			if _, hasProps := m["properties"]; hasProps {
				continue
			}
		}
		fallback = append(fallback, elem)
	}

	return map[string]any{"items": fallback}, nil
}
