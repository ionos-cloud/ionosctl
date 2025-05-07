package configgen

import (
	"encoding/json"
	"net"
	"net/http"
	"reflect"
	"testing"
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
		opts  FilterOptions
		expec []string
	}{
		{desc: "no filters returns all", opts: FilterOptions{}, expec: []string{"A", "B", "C"}},
		{desc: "version filter", opts: FilterOptions{Version: ptr("v1")}, expec: []string{"A", "C"}},
		{desc: "visibility filter", opts: FilterOptions{Visibility: ptr("public")}, expec: []string{"A", "C"}},
		{desc: "gate filter", opts: FilterOptions{Gate: ptr("Beta")}, expec: []string{"B"}},
		{desc: "whitelist filter", opts: FilterOptions{Whitelist: map[string]bool{"B": true}}, expec: []string{"B"}},
		{desc: "blacklist filter", opts: FilterOptions{Blacklist: map[string]bool{"C": true}}, expec: []string{"A", "B"}},
		{desc: "combined version and blacklist", opts: FilterOptions{Version: ptr("v1"), Blacklist: map[string]bool{"A": true}}, expec: []string{"C"}},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			out := filterPages(pages, tc.opts)
			var names []string
			for _, p := range out {
				names = append(names, p.Name)
			}
			if !reflect.DeepEqual(names, tc.expec) {
				t.Errorf("%s: expected %+v, got %+v", tc.desc, tc.expec, names)
			}
		})
	}
}

func TestToEndpoint(t *testing.T) {
	tests := []struct {
		srv   serverRaw
		expec Endpoint
	}{
		// relative URL yields default host
		{srv: serverRaw{URL: "/path"}, expec: Endpoint{Name: "api.ionos.com", SkipTLSVerify: false}},
		// invalid URL yields default
		{srv: serverRaw{URL: ":://bad"}, expec: Endpoint{Name: "api.ionos.com", SkipTLSVerify: false}},
		// absolute without region
		{srv: serverRaw{URL: "https://api.ionos.com"}, expec: Endpoint{Name: "api.ionos.com", SkipTLSVerify: false}},
		// with region
		{srv: serverRaw{URL: "https://nfs.de-fra.ionos.com"}, expec: Endpoint{Name: "nfs.de-fra.ionos.com", Location: "de/fra", SkipTLSVerify: false}},
		// multi-part host
		{srv: serverRaw{URL: "https://something.eu-wdc.ionos.com"}, expec: Endpoint{Name: "something.eu-wdc.ionos.com", Location: "eu/wdc", SkipTLSVerify: false}},
	}

	for _, tc := range tests {
		t.Run(tc.srv.URL, func(t *testing.T) {
			got := toEndpoint(tc.srv)
			if !reflect.DeepEqual(got, tc.expec) {
				t.Errorf("For URL %s expected %+v, got %+v", tc.srv.URL, tc.expec, got)
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

	// call GenerateConfig with no whitelist (includes both), but we whitelist only vpn
	opts := FilterOptions{Whitelist: map[string]bool{"vpn": true}, Visibility: ptr("public"), Gate: ptr("General-Availability"), Version: ptr("v1")}
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
          - location: de/fra
            name: nfs.de-fra.ionos.com
            skipTlsVerify: false
          - name: api.ionos.com
            skipTlsVerify: false
`

	if string(out) != expected {
		t.Errorf("unexpected YAML output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}
