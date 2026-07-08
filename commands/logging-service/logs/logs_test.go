package logs

import (
	"testing"

	logging "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
)

// TestFlattenPipelineLogs verifies that logs from multiple pipelines are
// flattened into a single {"items": [...]} payload, in order, with each log
// tagged by its parent pipeline ID.
func TestFlattenPipelineLogs(t *testing.T) {
	pipelines := logging.PipelineReadList{
		Items: []logging.PipelineRead{
			{Id: "p1", Properties: logging.Pipeline{Logs: []logging.PipelineNoAddrLogs{
				{Tag: "a", Source: "kubernetes"},
				{Tag: "b", Source: "docker"},
			}}},
			{Id: "p2", Properties: logging.Pipeline{Logs: []logging.PipelineNoAddrLogs{
				{Tag: "c", Source: "systemd"},
			}}},
		},
	}

	out := flattenPipelineLogs(pipelines)
	items, ok := out["items"].([]any)
	if !ok {
		t.Fatalf("items is not []any, got %T", out["items"])
	}
	if len(items) != 3 {
		t.Fatalf("want 3 logs, got %d", len(items))
	}

	want := []struct{ tag, pid string }{{"a", "p1"}, {"b", "p1"}, {"c", "p2"}}
	for i, w := range want {
		m, ok := items[i].(map[string]any)
		if !ok {
			t.Fatalf("item %d is not map[string]any, got %T", i, items[i])
		}
		if m["tag"] != w.tag {
			t.Errorf("item %d tag = %v, want %v", i, m["tag"], w.tag)
		}
		if m["_pipelineId"] != w.pid {
			t.Errorf("item %d _pipelineId = %v, want %v", i, m["_pipelineId"], w.pid)
		}
	}
}

// TestFlattenPipelineLogsEmpty verifies a non-nil empty items slice when there
// are no pipelines, so JSON output renders [] rather than null.
func TestFlattenPipelineLogsEmpty(t *testing.T) {
	out := flattenPipelineLogs(logging.PipelineReadList{})
	items, ok := out["items"].([]any)
	if !ok {
		t.Fatalf("items is not []any, got %T", out["items"])
	}
	if len(items) != 0 {
		t.Fatalf("want 0 logs, got %d", len(items))
	}
}
