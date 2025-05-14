package configgen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// indexURL is the source for the JSON index of OpenAPI specs
const indexURL = "https://ionos-cloud.github.io/rest-api/private-index.json"

// FilterOptions controls which APIs to include. Nil means "no filter".
type FilterOptions struct {
	Version     *string           // e.g. "v1"
	Visibility  *string           // e.g. "public"
	Gate        *string           // e.g. "General-Availability"
	Whitelist   map[string]bool   // API names to explicitly include
	Blacklist   map[string]bool   // API names to explicitly exclude
	CustomNames map[string]string // map spec-name -> desired name
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

	// build environment
	env := Environment{Name: "prod"}
	for _, page := range pages {
		// Construct full spec URL (indexURL base + page.Spec)
		base := strings.TrimSuffix(indexURL, "/rest-api/private-index.json")
		specURL := base + page.Spec

		// Load servers from spec
		servers, err := loadSpecServers(specURL)
		if err != nil {
			return nil, err
		}

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

func filterPages(pages []indexPage, opts FilterOptions) []indexPage {
	latest := make(map[string]indexPage)

	for _, p := range pages {
		name := p.Name
		if opts.CustomNames != nil {
			if custom, ok := opts.CustomNames[p.Name]; ok {
				name = custom
			}
		}

		if opts.Visibility != nil && p.Visibility != *opts.Visibility {
			continue
		}
		if opts.Gate != nil && p.Gate != *opts.Gate {
			continue
		}

		if opts.Whitelist != nil && !opts.Whitelist[name] {
			continue
		}
		if opts.Blacklist != nil && opts.Blacklist[name] {
			continue
		}

		if opts.Version != nil && p.Version != *opts.Version {
			continue
		}

		// collapse into latest[name]
		prev, exists := latest[name]
		if !exists || !compareVersions(p.Version, prev.Version) {
			p.Name = name
			latest[name] = p
		}
	}

	// collect results in a slice
	result := make([]indexPage, 0, len(latest))
	for _, p := range latest {
		result = append(result, p)
	}
	return result
}

func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		// fall back to string comparison
		if err1 != nil || err2 != nil {
			return parts1[i] < parts2[i]
		}

		if num1 != num2 {
			return num1 < num2
		}
	}

	// if all parts equal, the version with fewer parts is considered older
	return len(parts1) < len(parts2)
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

	u, err := url.Parse(s.URL)
	if err != nil {
		// If malformed, just return the raw URL as the name
		ep.Name = s.URL
		return ep
	}

	if !u.IsAbs() {
		// Relative URL (e.g., "/reseller/v2")
		ep.Name = "https://api.ionos.com" + s.URL
		return ep
	}

	ep.Name = u.String()
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) > 2 {
		ep.Location = strings.Join(parts[1:len(parts)-2], "/")
		ep.Location = strings.ReplaceAll(ep.Location, "-", "/")
	}

	return ep
}
