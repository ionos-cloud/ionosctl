package table

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// resetViper clears relevant viper keys for test isolation.
func resetViper(t *testing.T) {
	t.Helper()
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgQuiet, nil)
	viper.Set(constants.ArgNoHeaders, false)
	viper.Set(constants.FlagQuery, nil)
}

// --- Column helper tests ---

func TestAllCols(t *testing.T) {
	cols := []Column{
		{Name: "Id", Default: true},
		{Name: "Name", Default: true},
		{Name: "State", Default: false},
	}
	assert.Equal(t, []string{"Id", "Name", "State"}, AllCols(cols))
}

func TestDefaultCols(t *testing.T) {
	cols := []Column{
		{Name: "Id", Default: true},
		{Name: "Name", Default: true},
		{Name: "State", Default: false},
		{Name: "Description", Default: false},
	}
	assert.Equal(t, []string{"Id", "Name"}, DefaultCols(cols))
}

func TestResolveCols_Defaults(t *testing.T) {
	cols := []Column{
		{Name: "Id", Default: true},
		{Name: "Name", Default: true},
		{Name: "State"},
	}
	assert.Equal(t, []string{"Id", "Name"}, ResolveCols(cols, nil))
}

func TestResolveCols_All(t *testing.T) {
	cols := []Column{
		{Name: "Id", Default: true},
		{Name: "Name", Default: true},
		{Name: "State"},
	}
	assert.Equal(t, []string{"Id", "Name", "State"}, ResolveCols(cols, []string{"all"}))
}

func TestResolveCols_Custom(t *testing.T) {
	cols := []Column{
		{Name: "Id", Default: true},
		{Name: "Name", Default: true},
		{Name: "State"},
	}
	assert.Equal(t, []string{"State", "Name"}, ResolveCols(cols, []string{"State", "Name"}))
}

func TestResolveCols_CaseInsensitive(t *testing.T) {
	cols := []Column{
		{Name: "ServerId", Default: true},
		{Name: "Name", Default: true},
	}
	assert.Equal(t, []string{"ServerId"}, ResolveCols(cols, []string{"serverid"}))
}

func TestResolveCols_InvalidFallsBackToDefault(t *testing.T) {
	cols := []Column{
		{Name: "Id", Default: true},
		{Name: "Name", Default: true},
	}
	assert.Equal(t, []string{"Id", "Name"}, ResolveCols(cols, []string{"nonexistent"}))
}

func TestColsMessage(t *testing.T) {
	cols := []Column{
		{Name: "Id"},
		{Name: "Name"},
	}
	msg := ColsMessage(cols)
	assert.Contains(t, msg, "Id")
	assert.Contains(t, msg, "Name")
}

// --- Extract tests ---

var testCols = []Column{
	{Name: "DatacenterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Description", JSONPath: "properties.description"},
}

func singleDC() map[string]any {
	return map[string]any{
		"id":   "dc-123",
		"href": "/cloudapi/v6/datacenters/dc-123",
		"metadata": map[string]any{
			"state": "AVAILABLE",
		},
		"properties": map[string]any{
			"name":        "My DC",
			"location":    "de/txl",
			"description": "Test datacenter",
		},
	}
}

func listDCs() map[string]any {
	return map[string]any{
		"id":   "datacenters",
		"type": "collection",
		"items": []any{
			map[string]any{
				"id":       "dc-1",
				"metadata": map[string]any{"state": "AVAILABLE"},
				"properties": map[string]any{
					"name":     "DC One",
					"location": "de/txl",
				},
			},
			map[string]any{
				"id":       "dc-2",
				"metadata": map[string]any{"state": "BUSY"},
				"properties": map[string]any{
					"name":     "DC Two",
					"location": "us/las",
				},
			},
		},
	}
}

func TestExtract_SingleResource(t *testing.T) {
	tbl := New(testCols)
	err := tbl.Extract(singleDC())
	assert.NoError(t, err)
	assert.Len(t, tbl.Rows(), 1)

	row := tbl.Rows()[0]
	assert.Equal(t, "dc-123", row["DatacenterId"])
	assert.Equal(t, "My DC", row["Name"])
	assert.Equal(t, "de/txl", row["Location"])
	assert.Equal(t, "AVAILABLE", row["State"])
	assert.Equal(t, "Test datacenter", row["Description"])
}

func TestExtract_ListWithPrefix(t *testing.T) {
	tbl := New(testCols, WithPrefix("items"))
	err := tbl.Extract(listDCs())
	assert.NoError(t, err)
	assert.Len(t, tbl.Rows(), 2)

	assert.Equal(t, "dc-1", tbl.Rows()[0]["DatacenterId"])
	assert.Equal(t, "DC One", tbl.Rows()[0]["Name"])
	assert.Equal(t, "AVAILABLE", tbl.Rows()[0]["State"])

	assert.Equal(t, "dc-2", tbl.Rows()[1]["DatacenterId"])
	assert.Equal(t, "DC Two", tbl.Rows()[1]["Name"])
	assert.Equal(t, "BUSY", tbl.Rows()[1]["State"])
}

func TestExtract_Slice(t *testing.T) {
	data := []map[string]any{
		{"id": "dc-1", "properties": map[string]any{"name": "One"}, "metadata": map[string]any{"state": "AVAILABLE"}},
		{"id": "dc-2", "properties": map[string]any{"name": "Two"}, "metadata": map[string]any{"state": "BUSY"}},
	}
	tbl := New(testCols)
	err := tbl.Extract(data)
	assert.NoError(t, err)
	assert.Len(t, tbl.Rows(), 2)
	assert.Equal(t, "One", tbl.Rows()[0]["Name"])
	assert.Equal(t, "Two", tbl.Rows()[1]["Name"])
}

func TestExtract_EmptySlice(t *testing.T) {
	tbl := New(testCols)
	err := tbl.Extract([]map[string]any{})
	assert.NoError(t, err)
	assert.Empty(t, tbl.Rows())
}

func TestExtract_Raw(t *testing.T) {
	dc := singleDC()
	tbl := New(testCols)
	_ = tbl.Extract(dc)
	assert.Equal(t, dc, tbl.Raw())
}

func TestExtract_NilData(t *testing.T) {
	tbl := New(testCols)
	err := tbl.Extract(nil)
	assert.Error(t, err)
}

func TestExtract_WithFormatFunc(t *testing.T) {
	cols := []Column{
		{Name: "Id", JSONPath: "id", Default: true},
		{Name: "MaintenanceWindow", Default: true, Format: func(item map[string]any) any {
			mw := Navigate(item, "properties.maintenanceWindow")
			if mw == nil {
				return ""
			}
			m := mw.(map[string]any)
			return fmt.Sprintf("%s %s", m["dayOfTheWeek"], m["time"])
		}},
	}

	data := map[string]any{
		"id": "cluster-1",
		"properties": map[string]any{
			"maintenanceWindow": map[string]any{
				"dayOfTheWeek": "Monday",
				"time":         "15:00:00",
			},
		},
	}

	tbl := New(cols)
	err := tbl.Extract(data)
	assert.NoError(t, err)
	assert.Len(t, tbl.Rows(), 1)
	assert.Equal(t, "cluster-1", tbl.Rows()[0]["Id"])
	assert.Equal(t, "Monday 15:00:00", tbl.Rows()[0]["MaintenanceWindow"])
}

func TestExtract_FormatOverridesJSONPath(t *testing.T) {
	cols := []Column{
		{Name: "RAM", JSONPath: "properties.ram", Default: true, Format: func(item map[string]any) any {
			ram := Navigate(item, "properties.ram")
			if ram == nil {
				return ""
			}
			return fmt.Sprintf("%d GB", int64(ram.(float64))/1024/1024/1024)
		}},
	}

	data := map[string]any{
		"properties": map[string]any{
			"ram": float64(4294967296), // 4 GiB in bytes
		},
	}

	tbl := New(cols)
	_ = tbl.Extract(data)
	assert.Equal(t, "4 GB", tbl.Rows()[0]["RAM"])
}

func TestExtract_HrefAsParentColumn(t *testing.T) {
	cols := []Column{
		{Name: "ServerId", JSONPath: "id", Default: true},
		{Name: "DatacenterId", JSONPath: "href", Default: true},
	}

	data := map[string]any{
		"id":   "server-1",
		"href": "/cloudapi/v6/datacenters/dc-123/servers/server-1",
	}

	tbl := New(cols)
	_ = tbl.Extract(data)
	assert.Equal(t, "server-1", tbl.Rows()[0]["ServerId"])
	assert.Equal(t, "dc-123", tbl.Rows()[0]["DatacenterId"])
}

// --- SetCell tests ---

func TestSetCell(t *testing.T) {
	tbl := New(testCols)
	_ = tbl.Extract(singleDC())
	assert.Equal(t, "AVAILABLE", tbl.Rows()[0]["State"])

	tbl.SetCell(0, "State", "BUSY")
	assert.Equal(t, "BUSY", tbl.Rows()[0]["State"])
}

func TestSetCell_OutOfRange(t *testing.T) {
	tbl := New(testCols)
	_ = tbl.Extract(singleDC())
	// Should not panic
	tbl.SetCell(99, "State", "BUSY")
	tbl.SetCell(-1, "State", "BUSY")
	assert.Equal(t, "AVAILABLE", tbl.Rows()[0]["State"])
}

// --- Render tests ---

func TestRender_Text(t *testing.T) {
	resetViper(t)
	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId", "Name", "State"})
	assert.NoError(t, err)
	assert.Contains(t, out, "DatacenterId")
	assert.Contains(t, out, "dc-123")
	assert.Contains(t, out, "My DC")
	assert.Contains(t, out, "AVAILABLE")
}

func TestRender_TextNoHeaders(t *testing.T) {
	resetViper(t)
	viper.Set(constants.ArgNoHeaders, true)

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId", "Name"})
	assert.NoError(t, err)
	assert.NotContains(t, out, "DatacenterId") // Header should not appear
	assert.Contains(t, out, "dc-123")
	assert.Contains(t, out, "My DC")
}

func TestRender_TextEmptyColsEliminated(t *testing.T) {
	resetViper(t)
	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	// Request "Description" column which has a value, and a non-existent column
	out, err := tbl.Render([]string{"DatacenterId", "Description"})
	assert.NoError(t, err)
	assert.Contains(t, out, "dc-123")
	assert.Contains(t, out, "Test datacenter")
}

func TestRender_TextAllEmptyCols(t *testing.T) {
	resetViper(t)

	cols := []Column{
		{Name: "Missing", JSONPath: "nonexistent.field", Default: true},
	}
	data := map[string]any{"id": "1"}

	tbl := New(cols)
	_ = tbl.Extract(data)

	out, err := tbl.Render([]string{"Missing"})
	assert.NoError(t, err)
	assert.Empty(t, out)
}

func TestRender_APIJSON(t *testing.T) {
	resetViper(t)
	viper.Set(constants.ArgOutput, "api-json")

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId", "Name"})
	assert.NoError(t, err)
	assert.Contains(t, out, `"id": "dc-123"`)
	assert.Contains(t, out, `"name": "My DC"`)
}

func TestRender_JSON(t *testing.T) {
	resetViper(t)
	viper.Set(constants.ArgOutput, "json")

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId"})
	assert.NoError(t, err)
	// Single resource, should be passed through as-is
	assert.Contains(t, out, `"id": "dc-123"`)
}

func TestRender_Quiet(t *testing.T) {
	resetViper(t)
	viper.Set(constants.ArgQuiet, true)

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId"})
	assert.NoError(t, err)
	assert.Empty(t, out)
}

func TestRender_QuietFalseShowsOutput(t *testing.T) {
	resetViper(t)
	viper.Set(constants.ArgQuiet, false)

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId", "Name"})
	assert.NoError(t, err)
	assert.Contains(t, out, "dc-123")
	assert.Contains(t, out, "My DC")
}

func TestRender_TextWithEmptyQueryAllowed(t *testing.T) {
	resetViper(t)
	viper.Set(constants.FlagQuery, "")

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId", "Name"})
	assert.NoError(t, err)
	assert.Contains(t, out, "dc-123")
}

func TestRender_InvalidFormat(t *testing.T) {
	resetViper(t)
	viper.Set(constants.ArgOutput, "xml")

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	_, err := tbl.Render([]string{"DatacenterId"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}

func TestRender_TextListWithPrefix(t *testing.T) {
	resetViper(t)
	tbl := New(testCols, WithPrefix("items"))
	_ = tbl.Extract(listDCs())

	out, err := tbl.Render([]string{"DatacenterId", "Name", "State"})
	assert.NoError(t, err)
	assert.Contains(t, out, "dc-1")
	assert.Contains(t, out, "DC One")
	assert.Contains(t, out, "dc-2")
	assert.Contains(t, out, "DC Two")
	// Both rows should appear
	lines := strings.Split(strings.TrimSpace(out), "\n")
	assert.Equal(t, 3, len(lines)) // header + 2 rows
}

func TestRender_SliceFormatsCorrectly(t *testing.T) {
	resetViper(t)
	cols := []Column{
		{Name: "Name", JSONPath: "properties.name", Default: true},
		{Name: "Ips", JSONPath: "properties.ips", Default: true},
	}
	data := map[string]any{
		"properties": map[string]any{
			"name": "Test",
			"ips":  []any{"1.2.3.4", "5.6.7.8"},
		},
	}

	tbl := New(cols)
	_ = tbl.Extract(data)
	out, err := tbl.Render([]string{"Name", "Ips"})
	assert.NoError(t, err)
	assert.Contains(t, out, "1.2.3.4, 5.6.7.8")
}

func TestRender_FloatAsInt(t *testing.T) {
	resetViper(t)
	cols := []Column{
		{Name: "Cores", JSONPath: "properties.cores", Default: true},
	}
	data := map[string]any{
		"properties": map[string]any{
			"cores": float64(4),
		},
	}

	tbl := New(cols)
	_ = tbl.Extract(data)
	out, err := tbl.Render([]string{"Cores"})
	assert.NoError(t, err)
	assert.Contains(t, out, "4")
	assert.NotContains(t, out, "4.0")
}

// --- BeforeRender hook test ---

func TestRender_BeforeRenderHookSuppresses(t *testing.T) {
	resetViper(t)
	origHook := BeforeRender
	defer func() { BeforeRender = origHook }()

	BeforeRender = func(tbl *Table, visibleCols []string) bool {
		return false // suppress output
	}

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId"})
	assert.NoError(t, err)
	assert.Empty(t, out)
}

func TestRender_BeforeRenderHookAllows(t *testing.T) {
	resetViper(t)
	origHook := BeforeRender
	defer func() { BeforeRender = origHook }()

	BeforeRender = func(tbl *Table, visibleCols []string) bool {
		return true // allow output
	}

	tbl := New(testCols)
	_ = tbl.Extract(singleDC())

	out, err := tbl.Render([]string{"DatacenterId", "Name"})
	assert.NoError(t, err)
	assert.Contains(t, out, "dc-123")
}

// --- Sprint convenience test ---

func TestSprint(t *testing.T) {
	resetViper(t)
	out, err := Sprint(testCols, singleDC(), nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "dc-123")
	assert.Contains(t, out, "My DC")
}

func TestSprint_WithPrefix(t *testing.T) {
	resetViper(t)
	out, err := Sprint(testCols, listDCs(), nil, WithPrefix("items"))
	assert.NoError(t, err)
	assert.Contains(t, out, "DC One")
	assert.Contains(t, out, "DC Two")
}

func TestSprint_WithCustomCols(t *testing.T) {
	resetViper(t)
	out, err := Sprint(testCols, singleDC(), []string{"Name", "State"})
	assert.NoError(t, err)
	assert.Contains(t, out, "My DC")
	assert.Contains(t, out, "AVAILABLE")
	// DatacenterId should not be in output since we only requested Name and State
	assert.NotContains(t, out, "DatacenterId")
}

// --- Navigate helper test ---

func TestNavigate(t *testing.T) {
	m := map[string]any{
		"properties": map[string]any{
			"maintenanceWindow": map[string]any{
				"dayOfTheWeek": "Monday",
				"time":         "15:00",
			},
		},
	}
	assert.Equal(t, "Monday", Navigate(m, "properties.maintenanceWindow.dayOfTheWeek"))
	assert.Equal(t, "15:00", Navigate(m, "properties.maintenanceWindow.time"))
	assert.Nil(t, Navigate(m, "nonexistent.path"))
	assert.Nil(t, Navigate(m, "properties.nonexistent"))
}

func TestNavigate_Empty(t *testing.T) {
	assert.Nil(t, Navigate(nil, "a.b"))
	assert.Nil(t, Navigate(map[string]any{}, "a"))
}

// --- Internal helper tests ---

func TestParentIDFromHref(t *testing.T) {
	assert.Equal(t, "dc-123", parentIDFromHref("/cloudapi/v6/datacenters/dc-123/servers/srv-1"))
	assert.Equal(t, "", parentIDFromHref(nil))
	assert.Equal(t, "", parentIDFromHref(42))
	assert.Equal(t, "", parentIDFromHref("/short"))
}

func TestEliminateEmptyCols(t *testing.T) {
	rows := []map[string]any{
		{"A": "val", "B": nil, "C": ""},
		{"A": "val2", "B": nil, "C": "has-value"},
	}
	active := eliminateEmptyCols([]string{"A", "B", "C"}, rows)
	assert.Equal(t, []string{"A", "C"}, active)
}

func TestEliminateEmptyCols_EmptySlice(t *testing.T) {
	rows := []map[string]any{
		{"A": "val", "B": []any{}},
	}
	active := eliminateEmptyCols([]string{"A", "B"}, rows)
	assert.Equal(t, []string{"A"}, active)
}

func TestFormatCellValue(t *testing.T) {
	assert.Equal(t, "", formatCellValue(nil))
	assert.Equal(t, "hello", formatCellValue("hello"))
	assert.Equal(t, "42", formatCellValue(float64(42)))
	assert.Equal(t, "3.14", formatCellValue(3.14))
	assert.Equal(t, "true", formatCellValue(true))
	assert.Equal(t, "a, b, c", formatCellValue([]any{"a", "b", "c"}))
}

// --- Re-extract test (simulates --wait re-render) ---

func TestReExtract_UpdatesRows(t *testing.T) {
	resetViper(t)
	tbl := New(testCols)

	// First extract: BUSY state
	busyDC := map[string]any{
		"id":       "dc-123",
		"metadata": map[string]any{"state": "BUSY"},
		"properties": map[string]any{
			"name":     "My DC",
			"location": "de/txl",
		},
	}
	_ = tbl.Extract(busyDC)
	assert.Equal(t, "BUSY", tbl.Rows()[0]["State"])

	// Re-extract: AVAILABLE state (simulates --wait re-fetch)
	availableDC := map[string]any{
		"id":       "dc-123",
		"metadata": map[string]any{"state": "AVAILABLE"},
		"properties": map[string]any{
			"name":     "My DC",
			"location": "de/txl",
		},
	}
	_ = tbl.Extract(availableDC)
	assert.Equal(t, "AVAILABLE", tbl.Rows()[0]["State"])

	out, err := tbl.Render([]string{"DatacenterId", "Name", "State"})
	assert.NoError(t, err)
	assert.Contains(t, out, "AVAILABLE")
	assert.NotContains(t, out, "BUSY")
}

// --- Rerenderable interface test ---

func TestTable_ImplementsRerenderable(t *testing.T) {
	// Compile-time check that *Table satisfies the Rerenderable interface
	// (Extract + Render methods with the expected signatures)
	var _ interface {
		Extract(sourceData any) error
		Render(visibleCols []string) (string, error)
	} = (*Table)(nil)
}
