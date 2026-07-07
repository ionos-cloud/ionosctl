package core

import (
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
)

func TestFindRegionalConfig(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *cobra.Command
		wantLocs  []string
		wantURL   string
		wantProds []string
		wantFound bool
	}{
		{
			name: "annotations on direct command",
			setup: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				cmd.Annotations = map[string]string{
					AnnotationLocations:   "de/fra,de/txl,es/vit",
					AnnotationTemplateURL: "https://svc.%s.ionos.com",
				}
				return cmd
			},
			wantLocs:  []string{"de/fra", "de/txl", "es/vit"},
			wantURL:   "https://svc.%s.ionos.com",
			wantFound: true,
		},
		{
			name: "annotations on parent command",
			setup: func() *cobra.Command {
				parent := &cobra.Command{Use: "kafka"}
				parent.Annotations = map[string]string{
					AnnotationLocations:   "de/fra,us/las",
					AnnotationTemplateURL: "https://kafka.%s.ionos.com",
				}
				child := &cobra.Command{Use: "cluster"}
				leaf := &cobra.Command{Use: "list"}
				child.AddCommand(leaf)
				parent.AddCommand(child)
				return leaf
			},
			wantLocs:  []string{"de/fra", "us/las"},
			wantURL:   "https://kafka.%s.ionos.com",
			wantFound: true,
		},
		{
			name: "annotations include product names",
			setup: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				cmd.Annotations = map[string]string{
					AnnotationLocations:    "de/fra,de/txl",
					AnnotationTemplateURL:  "https://svc.%s.ionos.com",
					AnnotationProductNames: "cloud,compute",
				}
				return cmd
			},
			wantLocs:  []string{"de/fra", "de/txl"},
			wantURL:   "https://svc.%s.ionos.com",
			wantProds: []string{"cloud", "compute"},
			wantFound: true,
		},
		{
			name: "no annotations",
			setup: func() *cobra.Command {
				return &cobra.Command{Use: "test"}
			},
			wantFound: false,
		},
		{
			name: "partial annotations (missing URL)",
			setup: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				cmd.Annotations = map[string]string{
					AnnotationLocations: "de/fra",
				}
				return cmd
			},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup()
			locs, url, prods, found := findRegionalConfig(cmd)

			if found != tt.wantFound {
				t.Errorf("found = %v, want %v", found, tt.wantFound)
			}
			if !tt.wantFound {
				return
			}
			if url != tt.wantURL {
				t.Errorf("url = %q, want %q", url, tt.wantURL)
			}
			if strings.Join(locs, ",") != strings.Join(tt.wantLocs, ",") {
				t.Errorf("locs = %v, want %v", locs, tt.wantLocs)
			}
			if strings.Join(prods, ",") != strings.Join(tt.wantProds, ",") {
				t.Errorf("productNames = %v, want %v", prods, tt.wantProds)
			}
		})
	}
}

// TestFindOverridenURLPerLocation guards the regional list --all fix: an
// explicit --api-url override must win over the per-location template for
// every location queried.
func TestFindOverridenURLPerLocation(t *testing.T) {
	const tmpl = "https://svc.%s.ionos.com"

	t.Run("api-url override wins for all locations", func(t *testing.T) {
		cmd := &cobra.Command{Use: "list"}
		cmd.Flags().String(constants.ArgServerUrl, tmpl, "")
		if err := cmd.Flags().Set(constants.ArgServerUrl, "https://override.example.com"); err != nil {
			t.Fatal(err)
		}

		for _, loc := range []string{"de/fra", "de/txl"} {
			got := findOverridenURL(cmd, []string{"cloud"}, tmpl, loc)
			if got != "https://override.example.com" {
				t.Errorf("loc %s: url = %q, want override", loc, got)
			}
		}
	})

	t.Run("no override falls back to per-location template", func(t *testing.T) {
		cmd := &cobra.Command{Use: "list"}
		cmd.Flags().String(constants.ArgServerUrl, tmpl, "")

		got := findOverridenURL(cmd, []string{"cloud"}, tmpl, "de/txl")
		if got != "https://svc.de-txl.ionos.com" {
			t.Errorf("url = %q, want per-location template", got)
		}
	})
}

// TestRequireExplicitLocation guards the child-resource fix: single-resource
// operations on regional APIs must demand --location instead of silently
// defaulting to the first allowed location.
func TestRequireExplicitLocation(t *testing.T) {
	regionalCmd := func() *cobra.Command {
		cmd := &cobra.Command{Use: "get"}
		cmd.Annotations = map[string]string{
			AnnotationLocations:   "de/fra,de/txl",
			AnnotationTemplateURL: "https://svc.%s.ionos.com",
		}
		cmd.Flags().String(constants.FlagLocation, "de/fra", "")
		return cmd
	}

	t.Run("regional command without --location errors", func(t *testing.T) {
		if err := requireExplicitLocation(regionalCmd()); err == nil {
			t.Error("expected error when --location unset on regional command")
		}
	})

	t.Run("regional command with --location set passes", func(t *testing.T) {
		cmd := regionalCmd()
		if err := cmd.Flags().Set(constants.FlagLocation, "de/txl"); err != nil {
			t.Fatal(err)
		}
		if err := requireExplicitLocation(cmd); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("non-regional command passes", func(t *testing.T) {
		cmd := &cobra.Command{Use: "get"}
		if err := requireExplicitLocation(cmd); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestAnnotateLocation guards the api-json provenance fix: each per-location
// response must carry a "location" field without losing its original data.
func TestAnnotateLocation(t *testing.T) {
	t.Run("adds location to a collection object", func(t *testing.T) {
		type collection struct {
			Type  string   `json:"type"`
			Items []string `json:"items"`
		}
		got, err := annotateLocation(collection{Type: "collection", Items: []string{"a"}}, "de/txl")
		if err != nil {
			t.Fatal(err)
		}
		m, ok := got.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T", got)
		}
		if m["location"] != "de/txl" {
			t.Errorf("location = %v, want de/txl", m["location"])
		}
		if m["type"] != "collection" {
			t.Errorf("original fields lost: %v", m)
		}
		items, ok := m["items"].([]any)
		if !ok || len(items) != 1 || items[0] != "a" {
			t.Errorf("items not preserved: %v", m["items"])
		}
	})

	t.Run("empty collection is kept and annotated", func(t *testing.T) {
		got, err := annotateLocation(map[string]any{"items": []any{}}, "es/vit")
		if err != nil {
			t.Fatal(err)
		}
		m := got.(map[string]any)
		if m["location"] != "es/vit" {
			t.Errorf("location = %v, want es/vit", m["location"])
		}
		if items, ok := m["items"].([]any); !ok || len(items) != 0 {
			t.Errorf("empty items not preserved: %v", m["items"])
		}
	})

	t.Run("non-object response returned unmodified", func(t *testing.T) {
		got, err := annotateLocation([]int{1, 2, 3}, "de/fra")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := got.(map[string]any); ok {
			t.Errorf("expected non-object passthrough, got map: %v", got)
		}
	})
}

// TestLocationStampedItems guards the -o json merge fix: each item must gain a
// top-level "location" field while keeping its original fields.
func TestLocationStampedItems(t *testing.T) {
	t.Run("stamps each item with location", func(t *testing.T) {
		data := map[string]any{
			"items": []any{
				map[string]any{"id": "a", "properties": map[string]any{"name": "x"}},
				map[string]any{"id": "b"},
			},
		}
		items, err := locationStampedItems(data, "de/txl")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 2 {
			t.Fatalf("want 2 items, got %d", len(items))
		}
		for _, it := range items {
			m := it.(map[string]any)
			if m["location"] != "de/txl" {
				t.Errorf("item %v missing location", m)
			}
		}
		// original fields preserved
		if items[0].(map[string]any)["id"] != "a" {
			t.Errorf("original id lost: %v", items[0])
		}
	})

	t.Run("no items array returns nil", func(t *testing.T) {
		items, err := locationStampedItems(map[string]any{"type": "collection"}, "de/fra")
		if err != nil {
			t.Fatal(err)
		}
		if items != nil {
			t.Errorf("want nil, got %v", items)
		}
	})
}
