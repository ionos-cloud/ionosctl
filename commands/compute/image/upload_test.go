package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupFTP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		// short forms
		{"fra", "fra"},
		{"vit", "vit"},
		{"lhr", "lhr"},
		{"fkb", "fkb"},
		// API-style forms
		{"de/fra", "fra"},
		{"es/vit", "vit"},
		{"gb/lhr", "lhr"},
		{"de/fkb", "fkb"},
		// suffixed
		{"fra/2", "fra-2"},
		{"de/fra/2", "fra-2"},
		// unknown short - passthrough
		{"xyz", "xyz"},
		// unknown API-style - heuristic
		{"aa/xyz", "xyz"},
		{"aa/xyz/3", "xyz-3"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, lookupFTP(tt.input))
		})
	}
}

func TestDeduplicateLocations(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"no duplicates", []string{"fra", "vit"}, []string{"fra", "vit"}},
		{"short and API-style same location", []string{"vit", "es/vit"}, []string{"vit"}},
		{"API-style and short same location", []string{"es/vit", "vit"}, []string{"es/vit"}},
		{"mixed with duplicates", []string{"fra", "de/fra", "vit"}, []string{"fra", "vit"}},
		{"all unique", []string{"fra", "vit", "lhr"}, []string{"fra", "vit", "lhr"}},
		{"empty", []string{}, []string{}},
		{"single", []string{"fra"}, []string{"fra"}},
		{"suffixed not duplicate", []string{"fra", "fra/2"}, []string{"fra", "fra/2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, deduplicateLocations(tt.input))
		})
	}
}

func TestLookupAPI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		// short forms -> canonical
		{"fra", "de/fra"},
		{"vit", "es/vit"},
		{"lhr", "gb/lhr"},
		{"las", "us/las"},
		{"par", "fr/par"},
		{"mci", "us/mci"},
		// already canonical -> passthrough
		{"de/fra", "de/fra"},
		{"es/vit", "es/vit"},
		{"de/fra/2", "de/fra/2"},
		// unknown API-style -> passthrough
		{"xx/yyy", "xx/yyy"},
		// unknown short without country guess -> passthrough
		{"xyz", "xyz"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, lookupAPI(tt.input))
		})
	}
}
