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
