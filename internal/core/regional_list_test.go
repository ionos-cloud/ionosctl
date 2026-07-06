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
