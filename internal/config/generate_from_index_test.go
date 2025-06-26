package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConfigE2E(t *testing.T) {
	// prepare fake index and specs
	index := indexFile{
		Pages: []indexPage{
			{Name: "vpn", Version: "v1", Visibility: "public", Gate: "GA", Spec: "/foo.yaml"},
			{Name: "vpn", Version: "v2", Visibility: "public", Gate: "GA", Spec: "/foo-v2.yaml"},
			{Name: "db", Version: "v1", Visibility: "public", Gate: "GA", Spec: "/bar.json"},
		},
	}
	indexData, _ := json.Marshal(index)
	specYAML := `servers:
- url: https://first.example.com
`
	specJSON := `{"servers":[{"url":"/local"}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest-api/private-index.json":
			w.Write(indexData)
		case "/foo.yaml":
			w.Write([]byte(specYAML))
		case "/foo-v2.yaml":
			// v2: should pick v2 over v1
			w.Write([]byte(`servers:
- url: https://second.example.com
`))
		case "/bar.json":
			w.Write([]byte(specJSON))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	// override indexURL
	origIndexURL := indexURL
	indexURL = ts.URL + "/rest-api/private-index.json"
	defer func() { indexURL = origIndexURL }()

	settings := ProfileSettings{Version: 1.23, ProfileName: "me", Token: "tok", Environment: "dev"}
	// Only include "vpn" so db is filtered out
	opts := Filters{Whitelist: map[string]bool{"vpn": true}}

	cfg, err := NewFromIndex(settings, opts)
	assert.NoError(t, err)
	assert.Equal(t, fileconfiguration.Version(1.23), cfg.Version)
	prod := cfg.Environments[0].Products
	// Should have picked v2 endpoint
	assert.Len(t, prod, 1)
	endpoints := prod[0].Endpoints
	assert.Len(t, endpoints, 1)
	assert.Equal(t, "https://second.example.com", endpoints[0].Name)
}

func TestGenerateConfig_NoMatch(t *testing.T) {
	// serve an empty index
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"pages":[]}`))
	}))
	defer ts.Close()
	orig := indexURL
	indexURL = ts.URL
	defer func() { indexURL = orig }()

	_, err := NewFromIndex(ProfileSettings{Token: "x"}, Filters{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no APIs match given filters")
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.2", "1.10", true}, // numeric vs numeric
		{"1.10", "1.2", false},
		{"1.2.3", "1.2.3", false}, // equal
		{"1.2", "1.2.0", true},    // shorter < longer
		{"foo", "bar", false},     // lexicographic
		{"1.a", "1.b", true},      // lexicographic fallback
	}
	for _, tc := range cases {
		t.Run(tc.a+"<"+tc.b, func(t *testing.T) {
			assert.Equal(t, tc.want, compareVersions(tc.a, tc.b))
		})
	}
}

func TestToEndpoint(t *testing.T) {
	// absolute URL with multi-part host
	raw := serverRaw{URL: "https://foo.ab-cde.ionos.com"}
	ep := toEndpoint(raw)
	assert.Equal(t, "https://foo.ab-cde.ionos.com", ep.Name)
	assert.Equal(t, "ab/cde", ep.Location)

	// relative URL
	raw2 := serverRaw{URL: "/api"}
	ep2 := toEndpoint(raw2)
	assert.Equal(t, "https://api.ionos.com/api", ep2.Name)
	assert.Empty(t, ep2.Location)

	// invalid URL
	raw3 := serverRaw{URL: "://bad"}
	ep3 := toEndpoint(raw3)
	assert.Equal(t, "://bad", ep3.Name)
	assert.Empty(t, ep3.Location)
}

func TestLoadSpecServers_Error(t *testing.T) {
	// serve invalid YAML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not: [valid"))
	}))
	defer ts.Close()

	_, err := loadSpecServers(ts.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not parse spec YAML")
}
