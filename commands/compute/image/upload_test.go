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

func TestPickFTPCredentials(t *testing.T) {
	tests := []struct {
		name                                                   string
		cfgUser, cfgPass, envUser, envPass, profUser, profPass string
		wantUser, wantPass                                     string
		wantOK                                                 bool
	}{
		{
			name: "client basic creds win", cfgUser: "cfg", cfgPass: "cfgpw",
			envUser: "env", envPass: "envpw", wantUser: "cfg", wantPass: "cfgpw", wantOK: true,
		},
		{
			// token-only client: cfg user/pass empty -> fall back to env (the reported scenario).
			name: "env fallback when client has token only", envUser: "env", envPass: "envpw",
			wantUser: "env", wantPass: "envpw", wantOK: true,
		},
		{
			name: "profile fallback when env incomplete", envUser: "env", // envPass missing
			profUser: "prof", profPass: "profpw", wantUser: "prof", wantPass: "profpw", wantOK: true,
		},
		{
			// pure token-only, nothing else set: FTP cannot authenticate.
			name: "incomplete pair -> not ok", cfgUser: "onlyuser", wantOK: false,
		},
		{
			name: "all empty -> not ok", wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, p, ok := pickFTPCredentials(tt.cfgUser, tt.cfgPass, tt.envUser, tt.envPass, tt.profUser, tt.profPass)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantUser, u)
			assert.Equal(t, tt.wantPass, p)
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
