package configgen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

// indexURL is the source for the JSON index of OpenAPI specs
const indexURL = "https://ionos-cloud.github.io/rest-api/private-index.json"

// FilterOptions controls which APIs to include. Nil means "no filter".
type FilterOptions struct {
	Version    *string         // e.g. "v1"
	Visibility *string         // e.g. "public"
	Gate       *string         // e.g. "General-Availability"
	Whitelist  map[string]bool // API names to explicitly include
	Blacklist  map[string]bool // API names to explicitly exclude
}

// indexPage represents one entry in private-index.json
type indexPage struct {
	Name       string `json:"name"`
	Spec       string `json:"spec"`
	Visibility string `json:"visibility"`
	Version    string `json:"version"`
	Gate       string `json:"gate"`
}

// indexFile wraps the full JSON index
type indexFile struct {
	Pages []indexPage `json:"pages"`
}

// serverRaw matches the "servers" list in each OpenAPI spec
type serverRaw struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description,omitempty"`
}

// Config structure for YAML output
type Config struct {
	Version        string        `yaml:"version"`
	CurrentProfile string        `yaml:"currentProfile"`
	Profiles       []interface{} `yaml:"profiles"`
	Environments   []Environment `yaml:"environments"`
}

type Environment struct {
	Name     string    `yaml:"name"`
	Products []Product `yaml:"products"`
}

type Product struct {
	Name      string     `yaml:"name"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Location            string `yaml:"location,omitempty"`
	Name                string `yaml:"name"`
	SkipTLSVerify       bool   `yaml:"skipTlsVerify"`
	CertificateAuthData string `yaml:"certificateAuthData,omitempty"`
}

// GenerateConfig builds the endpoints.yaml content based on the index and OpenAPI specs.
func GenerateConfig(opts FilterOptions) ([]byte, error) {
	// 1. Load and parse the index JSON
	idx, err := loadIndex()
	if err != nil {
		return nil, err
	}

	// 2. Filter pages
	pages := filterPages(idx.Pages, opts)
	if len(pages) == 0 {
		return nil, fmt.Errorf("no APIs match given filters")
	}

	fmt.Println("got", len(pages), "pages from index")

	// build environment
	env := Environment{Name: "prod"}
	for _, page := range pages {
		// Construct full spec URL (indexURL base + page.Spec)
		base := strings.TrimSuffix(indexURL, "/rest-api/private-index.json")
		specURL := base + page.Spec

		fmt.Println("loading spec", specURL)

		// Load servers from spec
		servers, err := loadSpecServers(specURL)
		if err != nil {
			return nil, err
		}

		fmt.Println("loading servers", servers)

		// Convert servers into endpoints
		prod := Product{Name: page.Name}
		for _, srv := range servers {
			ep := toEndpoint(srv)
			prod.Endpoints = append(prod.Endpoints, ep)
		}
		env.Products = append(env.Products, prod)
	}

	// assemble config
	cfg := Config{
		Version:        "1.0",
		CurrentProfile: "",
		Profiles:       []interface{}{},
		Environments:   []Environment{env},
	}

	var out strings.Builder
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(cfg); err != nil {
		return nil, fmt.Errorf("could not encode YAML: %w", err)
	}

	return []byte(out.String()), nil
}

// filterPages applies the filtering options to the list of pages
func filterPages(pages []indexPage, opts FilterOptions) []indexPage {
	var result []indexPage
	for _, p := range pages {
		if opts.Version != nil && p.Version != *opts.Version {
			continue
		}
		if opts.Visibility != nil && p.Visibility != *opts.Visibility {
			continue
		}
		if opts.Gate != nil && p.Gate != *opts.Gate {
			continue
		}
		if opts.Whitelist != nil && !opts.Whitelist[p.Name] {
			continue
		}
		if opts.Blacklist != nil && opts.Blacklist[p.Name] {
			continue
		}
		result = append(result, p)
	}
	return result
}

func loadIndex() (*indexFile, error) {
	resp, err := http.Get(indexURL)
	if err != nil {
		return nil, fmt.Errorf("fetch index: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read index body: %w", err)
	}

	var idx indexFile
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse index JSON: %w", err)
	}
	return &idx, nil
}

// loadSpecServers fetches an OpenAPI spec and returns its servers list
func loadSpecServers(urlStr string) ([]serverRaw, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("fetch spec: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read spec body: %w", err)
	}

	var wrapper struct {
		Servers []serverRaw `yaml:"servers"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parse spec YAML: %w", err)
	}
	return wrapper.Servers, nil
}

// toEndpoint converts a serverRaw into our Endpoint type
func toEndpoint(s serverRaw) Endpoint {
	ep := Endpoint{SkipTLSVerify: false}

	// parse URL; relative means default API
	u, err := url.Parse(s.URL)
	if err != nil || !u.IsAbs() {
		ep.Name = "api.ionos.com"
		return ep
	}

	ep.Name = u.Host

	// try extracting region: host like "nfs.de-fra.ionos.com"
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) >= 3 {
		region := parts[1] // e.g. "de-fra"
		sub := strings.Split(region, "-")
		if len(sub) == 2 {
			ep.Location = sub[0] + "/" + sub[1]
		}
	}

	return ep
}
