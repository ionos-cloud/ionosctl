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
	// Only include "vpn" so db is filtered out. Both v1 and v2 share the same
	// renamed name "vpn", so dedup keeps the latest (v2).
	opts := Filters{Whitelist: map[string]bool{"vpn": true}}

	cfg, err := NewFromIndex(settings, opts)
	assert.NoError(t, err)
	assert.Equal(t, fileconfiguration.Version(1.23), cfg.Version)
	prod := cfg.Environments[0].Products
	assert.Len(t, prod, 1)
	endpoints := prod[0].Endpoints
	assert.Len(t, endpoints, 1)
	assert.Equal(t, "https://second.example.com", endpoints[0].Name)
}

func TestGenerateConfigE2E_VersionedCustomNames(t *testing.T) {
	// When CustomNames maps different versions to different names, both should appear
	index := indexFile{
		Pages: []indexPage{
			{Name: "postgresql", Version: "v1", Visibility: "public", Gate: "GA", Spec: "/psql-v1.yaml"},
			{Name: "postgresql", Version: "v2", Visibility: "public", Gate: "EA", Spec: "/psql-v2.yaml"},
		},
	}
	indexData, _ := json.Marshal(index)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest-api/private-index.json":
			w.Write(indexData)
		case "/psql-v1.yaml":
			w.Write([]byte(`servers:
- url: https://psql-v1.example.com
`))
		case "/psql-v2.yaml":
			w.Write([]byte(`servers:
- url: https://psql-v2.example.com
`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	origIndexURL := indexURL
	indexURL = ts.URL + "/rest-api/private-index.json"
	defer func() { indexURL = origIndexURL }()

	settings := ProfileSettings{Token: "tok"}
	opts := Filters{
		CustomNames: map[string]string{
			"postgresql:v1": "psql",
			"postgresql:v2": "psqlv2",
		},
	}

	cfg, err := NewFromIndex(settings, opts)
	assert.NoError(t, err)
	prods := cfg.Environments[0].Products
	assert.Len(t, prods, 2)

	prodMap := make(map[string]string)
	for _, p := range prods {
		prodMap[p.Name] = p.Endpoints[0].Name
	}
	assert.Equal(t, "https://psql-v1.example.com", prodMap["psql"])
	assert.Equal(t, "https://psql-v2.example.com", prodMap["psqlv2"])
}

func TestFilterPages_VersionedWhitelist(t *testing.T) {
	pages := []indexPage{
		{Name: "postgresql", Version: "v1", Visibility: "public"},
		{Name: "postgresql", Version: "v2", Visibility: "public"},
		{Name: "vpn", Version: "v1", Visibility: "public"},
	}

	// whitelist only postgresql:v1 — should exclude v2 and vpn
	result := filterPages(pages, Filters{
		Whitelist: map[string]bool{"postgresql:v1": true},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "postgresql", result[0].Name)
	assert.Equal(t, "v1", result[0].Version)
}

func TestFilterPages_VersionedBlacklist(t *testing.T) {
	pages := []indexPage{
		{Name: "postgresql", Version: "v1", Visibility: "public"},
		{Name: "postgresql", Version: "v2", Visibility: "public"},
	}

	// blacklist only v1 — should keep v2
	result := filterPages(pages, Filters{
		Blacklist: map[string]bool{"postgresql:v1": true},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "v2", result[0].Version)
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
