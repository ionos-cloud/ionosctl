package core

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFindRegionalConfig(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *cobra.Command
		wantLocs   []string
		wantURL    string
		wantFound  bool
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
			locs, url, found := findRegionalConfig(cmd)

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
		})
	}
}
