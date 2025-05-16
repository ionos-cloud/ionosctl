package configgen

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptr(s string) *string { return &s }

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFilterPages(t *testing.T) {
	pages := []indexPage{
		{Name: "A", Version: "v1", Visibility: "public", Gate: "GA"},
		{Name: "B", Version: "v2", Visibility: "private", Gate: "Beta"},
		{Name: "C", Version: "v1", Visibility: "public", Gate: "GA"},
	}

	tests := []struct {
		desc  string
		opts  Filters
		expec []string
	}{
		{desc: "no filters returns all", opts: Filters{}, expec: []string{"A", "B", "C"}},
		{desc: "version filter", opts: Filters{Version: ptr("v1")}, expec: []string{"A", "C"}},
		{desc: "visibility filter", opts: Filters{Visibility: ptr("public")}, expec: []string{"A", "C"}},
		{desc: "gate filter", opts: Filters{Gate: ptr("Beta")}, expec: []string{"B"}},
		{desc: "whitelist filter", opts: Filters{Whitelist: map[string]bool{"B": true}}, expec: []string{"B"}},
		{desc: "blacklist filter", opts: Filters{Blacklist: map[string]bool{"C": true}}, expec: []string{"A", "B"}},
		{desc: "combined version and blacklist", opts: Filters{Version: ptr("v1"), Blacklist: map[string]bool{"A": true}}, expec: []string{"C"}},
		{desc: "custom names", opts: Filters{CustomNames: map[string]string{"A": "Alpha", "B": "Beta"}}, expec: []string{"Alpha", "Beta", "C"}},
		{desc: "custom names with version filter", opts: Filters{Version: ptr("v1"), CustomNames: map[string]string{"A": "Alpha", "B": "Beta"}}, expec: []string{"Alpha", "C"}},
		{desc: "whitelist the custom names, not original name", opts: Filters{Version: ptr("v1"), Whitelist: map[string]bool{"NewName": true}, CustomNames: map[string]string{"A": "NewName"}}, expec: []string{"NewName"}},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			out := filterPages(pages, tc.opts)
			var names []string
			for _, p := range out {
				names = append(names, p.Name)
			}
			if !assert.ElementsMatch(t, tc.expec, names) {
				t.Errorf("%s: expected %+v, got %+v", tc.desc, tc.expec, names)
			}
		})
	}
}

// TestGenerateConfigE2E spins up an HTTP client that serves a minimal index
// and corresponding spec files, then runs GenerateConfig and logs the YAML.
func TestGenerateConfigE2E(t *testing.T) {
	// prepare fake index and specs
	index := indexFile{
		Pages: []indexPage{
			{Name: "vpn", Version: "v1", Visibility: "public", Gate: "General-Availability", Spec: "/rest-api/foo.yaml"},
			{Name: "db", Version: "v1", Visibility: "public", Gate: "General-Availability", Spec: "/rest-api/bar.json"},
		},
	}
	indexData, _ := json.Marshal(index)

	// minimal OpenAPI spec YAML with servers list
	specYaml := []byte(`servers:
- url: https://foo.ab-cde.ionos.com
  description: AB/CDE location
- url: https://bar.vw-xyz.ionos.com
  description: VW/XYZ endpoint
`)

	specJson := []byte(`{"servers":[{"url":"/local"}]}`)

	// httptest server
	hs := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest-api/private-index.json":
			w.Write(indexData)
		case "/rest-api/foo.yaml":
			w.Write(specYaml)
		case "/rest-api/bar.json":
			w.Write(specJson)
		default:
			http.NotFound(w, r)
		}
	})
	ts := &http.Server{Addr: "127.0.0.1:0", Handler: hs}
	ln, err := net.Listen("tcp", ts.Addr)
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	go ts.Serve(ln)
	defer ts.Close()

	// override default transport to redirect indexURL & specs
	origTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		// rewrite host
		req2 := new(http.Request)
		*req2 = *req
		req2.URL.Scheme = "http"
		req2.URL.Host = ln.Addr().String()
		return origTransport.RoundTrip(req2)
	})
	defer func() { http.DefaultTransport = origTransport }()

	opts := Filters{Whitelist: map[string]bool{"vpn": true}, Visibility: ptr("public"), Gate: ptr("General-Availability"), Version: ptr("v1")}
	out, err := GenerateConfig(opts)
	if err != nil {
		t.Fatalf("GenerateConfig failed: %v", err)
	}

	t.Logf("Generated YAML:\n%s", string(out))

	expected := `version: "1.0"
currentProfile: ""
profiles: []
environments:
  - name: prod
    products:
      - name: vpn
        endpoints:
          - location: ab/cde
            name: https://foo.ab-cde.ionos.com
            skipTlsVerify: false
          - location: vw/xyz
            name: https://bar.vw-xyz.ionos.com
            skipTlsVerify: false
`
	assert.Equal(t, expected, string(out))

	out, err = GenerateConfig(Filters{})
	if err != nil {
		t.Fatalf("GenerateConfig failed: %v", err)
	}

	t.Logf("Generated YAML:\n%s", string(out))

	expected = `version: "1.0"
currentProfile: ""
profiles: []
environments:
  - name: prod
    products:
      - name: vpn
        endpoints:
          - location: ab/cde
            name: https://foo.ab-cde.ionos.com
            skipTlsVerify: false
          - location: vw/xyz
            name: https://bar.vw-xyz.ionos.com
            skipTlsVerify: false
      - name: db
        endpoints:
          - name: https://api.ionos.com/local
            skipTlsVerify: false
`
	assert.Equal(t, expected, string(out))

}
