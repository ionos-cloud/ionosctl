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

// TestFillOutEmptyFields is a regression test for the 'logs update' no-op bug:
// fields the user set must survive, and only the fields left unset may fall back
// to the existing log. Previously every field was overwritten with the old value,
// making --new-log-tag/--log-source/--log-protocol/--log-type/--log-retention-time
// no-ops.
func TestFillOutEmptyFields(t *testing.T) {
	oldLog := logging.PipelineNoAddrLogs{
		Tag:      "k8s",
		Source:   "kubernetes",
		Protocol: "http",
		Destinations: []logging.PipelineNoAddrLogsDestinations{
			{Type: "loki", RetentionInDays: 7},
		},
	}

	t.Run("renames tag, preserves the rest", func(t *testing.T) {
		// generatePatchObject always sets Destinations, so an unchanged
		// destination arrives as an empty entry (Type "", RetentionInDays 0).
		newLog := logging.PipelineNoAddrLogs{
			Tag:          "k8s-prod",
			Destinations: []logging.PipelineNoAddrLogsDestinations{{}},
		}

		got := fillOutEmptyFields(oldLog, newLog)

		if got.Tag != "k8s-prod" {
			t.Errorf("Tag = %q, want %q (user value must survive)", got.Tag, "k8s-prod")
		}
		if got.Source != "kubernetes" {
			t.Errorf("Source = %q, want %q (unset falls back to old)", got.Source, "kubernetes")
		}
		if got.Protocol != "http" {
			t.Errorf("Protocol = %q, want %q (unset falls back to old)", got.Protocol, "http")
		}
		if got.Destinations[0].Type != "loki" {
			t.Errorf("Destinations[0].Type = %q, want %q (unset falls back to old)", got.Destinations[0].Type, "loki")
		}
		if got.Destinations[0].RetentionInDays != 7 {
			t.Errorf("Destinations[0].RetentionInDays = %d, want 7 (unset falls back to old)", got.Destinations[0].RetentionInDays)
		}
	})

	t.Run("changes only retention", func(t *testing.T) {
		newLog := logging.PipelineNoAddrLogs{
			Destinations: []logging.PipelineNoAddrLogsDestinations{{RetentionInDays: 30}},
		}

		got := fillOutEmptyFields(oldLog, newLog)

		if got.Destinations[0].RetentionInDays != 30 {
			t.Errorf("RetentionInDays = %d, want 30 (user value must survive)", got.Destinations[0].RetentionInDays)
		}
		if got.Destinations[0].Type != "loki" {
			t.Errorf("Type = %q, want %q (unset falls back to old)", got.Destinations[0].Type, "loki")
		}
		if got.Tag != "k8s" || got.Source != "kubernetes" || got.Protocol != "http" {
			t.Errorf("unset identity fields changed: %+v", got)
		}
	})

	t.Run("empty update preserves the old log", func(t *testing.T) {
		newLog := logging.PipelineNoAddrLogs{
			Destinations: []logging.PipelineNoAddrLogsDestinations{{}},
		}

		got := fillOutEmptyFields(oldLog, newLog)

		if got.Tag != oldLog.Tag || got.Source != oldLog.Source || got.Protocol != oldLog.Protocol {
			t.Errorf("identity fields changed: got %+v, want %+v", got, oldLog)
		}
		if got.Destinations[0].Type != "loki" || got.Destinations[0].RetentionInDays != 7 {
			t.Errorf("destination changed: got %+v", got.Destinations[0])
		}
	})
}
