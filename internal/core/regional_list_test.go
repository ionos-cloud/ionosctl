package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TestMain initializes the client singleton with a dummy token so that
// client.Must() (reached via client.NewRegionalConfig inside ListAllLocations
// and RunForAllLocations) does not fatally exit during these unit tests. The
// token is never used for a real request - every test supplies a fake fetchFn.
func TestMain(m *testing.M) {
	os.Setenv(constants.EnvToken, "unit-test-token")
	// Force the once.Do init to run now, with the token present, so getClientErr
	// stays nil and Must() returns a usable client.
	_ = client.Must(func(error) {})
	os.Exit(m.Run())
}

// newRegionalTestCmd builds a CommandConfig whose cobra command carries the
// regional annotations, with output/error captured in the returned buffers.
func newRegionalTestCmd(locations []string) (*CommandConfig, *bytes.Buffer, *bytes.Buffer) {
	cmd := &cobra.Command{Use: "list"}
	cmd.Annotations = map[string]string{
		AnnotationLocations:   strings.Join(locations, ","),
		AnnotationTemplateURL: "https://svc.%s.ionos.com",
	}
	cmd.Flags().String(constants.FlagLocation, locations[0], "")
	cmd.Flags().StringSlice(constants.ArgCols, nil, "")
	var out, errb bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errb)
	return &CommandConfig{Command: &Command{Command: cmd}}, &out, &errb
}

// okItems is a fake fetchFn returning a collection with a single item.
func okItems(*shared.Configuration) (any, error) {
	return map[string]any{"items": []any{map[string]any{"id": "x"}}, "type": "collection"}, nil
}

func setOutputFormat(t *testing.T, format string) {
	t.Helper()
	prev := viper.GetString(constants.ArgOutput)
	viper.Set(constants.ArgOutput, format)
	t.Cleanup(func() { viper.Set(constants.ArgOutput, prev) })
}

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

// TestListAllLocations_TargetSelection verifies how many locations are queried:
// all when --location is unset, exactly one when it is set.
func TestListAllLocations_TargetSelection(t *testing.T) {
	setOutputFormat(t, "json")

	t.Run("queries all locations when --location unset", func(t *testing.T) {
		cc, _, _ := newRegionalTestCmd([]string{"de/fra", "de/txl", "es/vit"})
		var calls atomic.Int32
		err := cc.ListAllLocations(nil, func(cfg *shared.Configuration) (any, error) {
			calls.Add(1)
			return okItems(cfg)
		})
		if err != nil {
			t.Fatal(err)
		}
		if calls.Load() != 3 {
			t.Errorf("fetchFn called %d times, want 3", calls.Load())
		}
	})

	t.Run("queries one location when --location set", func(t *testing.T) {
		cc, _, _ := newRegionalTestCmd([]string{"de/fra", "de/txl", "es/vit"})
		if err := cc.Command.Command.Flags().Set(constants.FlagLocation, "de/txl"); err != nil {
			t.Fatal(err)
		}
		var calls atomic.Int32
		err := cc.ListAllLocations(nil, func(cfg *shared.Configuration) (any, error) {
			calls.Add(1)
			return okItems(cfg)
		})
		if err != nil {
			t.Fatal(err)
		}
		if calls.Load() != 1 {
			t.Errorf("fetchFn called %d times, want 1", calls.Load())
		}
	})
}

// TestListAllLocations_JSONStamping verifies merged -o json output stamps each
// item with its source location, and that the shape is the same whether or not
// --location is set (only the item count differs).
func TestListAllLocations_JSONStamping(t *testing.T) {
	setOutputFormat(t, "json")

	parseItems := func(t *testing.T, out string) []any {
		t.Helper()
		var payload map[string]any
		if err := json.Unmarshal([]byte(out), &payload); err != nil {
			t.Fatalf("output is not a JSON object: %v\n%s", err, out)
		}
		items, ok := payload["items"].([]any)
		if !ok {
			t.Fatalf("output has no items array: %s", out)
		}
		return items
	}

	t.Run("all locations, each item stamped", func(t *testing.T) {
		cc, out, _ := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		if err := cc.ListAllLocations(nil, okItems); err != nil {
			t.Fatal(err)
		}
		items := parseItems(t, out.String())
		if len(items) != 2 {
			t.Fatalf("want 2 merged items, got %d: %s", len(items), out.String())
		}
		got := map[string]bool{}
		for _, it := range items {
			got[it.(map[string]any)["location"].(string)] = true
		}
		if !got["de/fra"] || !got["de/txl"] {
			t.Errorf("items not stamped with both locations: %v", got)
		}
	})

	t.Run("single location has same shape", func(t *testing.T) {
		cc, out, _ := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		if err := cc.Command.Command.Flags().Set(constants.FlagLocation, "de/txl"); err != nil {
			t.Fatal(err)
		}
		if err := cc.ListAllLocations(nil, okItems); err != nil {
			t.Fatal(err)
		}
		items := parseItems(t, out.String())
		if len(items) != 1 {
			t.Fatalf("want 1 item, got %d", len(items))
		}
		if items[0].(map[string]any)["location"] != "de/txl" {
			t.Errorf("item not stamped with de/txl: %v", items[0])
		}
	})
}

// TestListAllLocations_APIJSON verifies -o api-json is an array with one
// annotated element per location.
func TestListAllLocations_APIJSON(t *testing.T) {
	setOutputFormat(t, "api-json")

	cc, out, _ := newRegionalTestCmd([]string{"de/fra", "de/txl"})
	if err := cc.ListAllLocations(nil, okItems); err != nil {
		t.Fatal(err)
	}

	var arr []map[string]any
	if err := json.Unmarshal(out.Bytes(), &arr); err != nil {
		t.Fatalf("api-json is not an array of objects: %v\n%s", err, out.String())
	}
	if len(arr) != 2 {
		t.Fatalf("want 2 per-location responses, got %d", len(arr))
	}
	got := map[string]bool{}
	for _, resp := range arr {
		loc, _ := resp["location"].(string)
		got[loc] = true
		if _, ok := resp["items"]; !ok {
			t.Errorf("per-location response lost its body: %v", resp)
		}
	}
	if !got["de/fra"] || !got["de/txl"] {
		t.Errorf("responses not annotated with both locations: %v", got)
	}
}

// TestListAllLocations_Errors covers error aggregation: a single target returns
// its error verbatim, while all-locations failures are wrapped.
func TestListAllLocations_Errors(t *testing.T) {
	setOutputFormat(t, "json")
	boom := func(*shared.Configuration) (any, error) { return nil, errors.New("boom") }

	t.Run("single target returns verbatim error", func(t *testing.T) {
		cc, _, _ := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		if err := cc.Command.Command.Flags().Set(constants.FlagLocation, "de/txl"); err != nil {
			t.Fatal(err)
		}
		err := cc.ListAllLocations(nil, boom)
		if err == nil || err.Error() != "boom" {
			t.Errorf("want verbatim \"boom\", got %v", err)
		}
	})

	t.Run("all locations failing is wrapped", func(t *testing.T) {
		cc, _, _ := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		err := cc.ListAllLocations(nil, boom)
		if err == nil || !strings.Contains(err.Error(), "failed to list from all locations") {
			t.Errorf("want wrapped all-locations error, got %v", err)
		}
	})
}

// TestListAllLocations_PartialFailure verifies a failing location warns on
// stderr while succeeding locations still render.
func TestListAllLocations_PartialFailure(t *testing.T) {
	setOutputFormat(t, "json")

	cc, out, errb := newRegionalTestCmd([]string{"de/fra", "de/txl"})
	err := cc.ListAllLocations(nil, func(cfg *shared.Configuration) (any, error) {
		if strings.Contains(cfg.Servers[0].URL, "de-txl") {
			return nil, errors.New("txl-down")
		}
		return okItems(cfg)
	})
	if err != nil {
		t.Fatalf("partial failure should not error out: %v", err)
	}
	if !strings.Contains(errb.String(), "de/txl") || !strings.Contains(errb.String(), "txl-down") {
		t.Errorf("expected warning for de/txl, stderr: %s", errb.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("bad output: %v", err)
	}
	if items := payload["items"].([]any); len(items) != 1 {
		t.Errorf("want 1 surviving item, got %d", len(items))
	}
}

// TestRunForAllLocations covers the delete --all fan-out helper: it runs once
// per location, aggregates errors, and announces the blast radius.
func TestRunForAllLocations(t *testing.T) {
	t.Run("runs once per location and announces blast radius", func(t *testing.T) {
		cc, _, errb := newRegionalTestCmd([]string{"de/fra", "de/txl", "es/vit"})
		var got []string
		err := cc.RunForAllLocations(func(_ *shared.Configuration, loc string) error {
			got = append(got, loc)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Join(got, ",") != "de/fra,de/txl,es/vit" {
			t.Errorf("ran for %v, want all three in order", got)
		}
		if !strings.Contains(errb.String(), "Operating across all 3 locations") {
			t.Errorf("missing blast-radius notice, stderr: %s", errb.String())
		}
	})

	t.Run("aggregates per-location errors", func(t *testing.T) {
		cc, _, _ := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		err := cc.RunForAllLocations(func(_ *shared.Configuration, loc string) error {
			if loc == "de/txl" {
				return errors.New("nope")
			}
			return nil
		})
		if err == nil || !strings.Contains(err.Error(), "location de/txl: nope") {
			t.Errorf("want aggregated error naming de/txl, got %v", err)
		}
	})

	t.Run("single target runs once without blast-radius notice", func(t *testing.T) {
		cc, _, errb := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		if err := cc.Command.Command.Flags().Set(constants.FlagLocation, "de/txl"); err != nil {
			t.Fatal(err)
		}
		var calls int
		var seen string
		err := cc.RunForAllLocations(func(_ *shared.Configuration, loc string) error {
			calls++
			seen = loc
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if calls != 1 || seen != "de/txl" {
			t.Errorf("want single run for de/txl, got calls=%d loc=%q", calls, seen)
		}
		if strings.Contains(errb.String(), "Operating across") {
			t.Errorf("single target should not announce blast radius, stderr: %s", errb.String())
		}
	})

	t.Run("blast-radius notice suppressed with --quiet", func(t *testing.T) {
		prev := viper.GetBool(constants.ArgQuiet)
		viper.Set(constants.ArgQuiet, true)
		t.Cleanup(func() { viper.Set(constants.ArgQuiet, prev) })

		cc, _, errb := newRegionalTestCmd([]string{"de/fra", "de/txl"})
		err := cc.RunForAllLocations(func(_ *shared.Configuration, _ string) error { return nil })
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(errb.String(), "Operating across") {
			t.Errorf("--quiet should suppress blast-radius notice, stderr: %s", errb.String())
		}
	})
}
