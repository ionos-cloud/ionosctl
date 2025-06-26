// cfggen.go
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"gopkg.in/yaml.v3"

	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
)

// indexURL is the source for the JSON index of OpenAPI specs
var indexURL = "https://ionos-cloud.github.io/rest-api/private-index.json"

// Filters controls which APIs to include
type Filters struct {
	Version    *string         // e.g. "v1"
	Visibility *string         // e.g. "public"
	Gate       *string         // e.g. "General-Availability"
	Whitelist  map[string]bool // API names to explicitly include
	Blacklist  map[string]bool // API names to explicitly exclude

	CustomNames map[string]string // map spec-name -> desired name
}

// ProfileSettings holds options for config generation
type ProfileSettings struct {
	Version     float64 // default: 1.0
	ProfileName string  // default: "user"
	Token       string  // default: "<token>"
	Environment string  // default: "prod"
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

// NewFromIndex builds a FileConfig based on the index and OpenAPI specs.
func NewFromIndex(settings ProfileSettings, opts Filters) (*fileconfiguration.FileConfig, error) {
	// default version
	if settings.Version == 0 {
		settings.Version = 1.0
	}
	// default token/profile/env
	if settings.Token == "" {
		settings.Token = "<token>"
	}
	if settings.ProfileName == "" {
		settings.ProfileName = "user"
	}
	if settings.Environment == "" {
		settings.Environment = "prod"
	}

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

	// Build environment products/endpoints
	envProducts := make([]fileconfiguration.Product, 0, len(pages))
	for _, page := range pages {
		base := strings.TrimSuffix(indexURL, "/rest-api/private-index.json")
		specURL := base + page.Spec

		servers, err := loadSpecServers(specURL)
		if err != nil {
			return nil, err
		}

		prod := fileconfiguration.Product{Name: page.Name}
		for _, srv := range servers {
			prod.Endpoints = append(prod.Endpoints, toEndpoint(srv))
		}
		envProducts = append(envProducts, prod)
	}

	// Sort products by name
	sort.Slice(envProducts, func(i, j int) bool {
		return envProducts[i].Name < envProducts[j].Name
	})

	// Assemble FileConfig
	fc := &fileconfiguration.FileConfig{
		Version:        settings.Version,
		CurrentProfile: settings.ProfileName,
		Profiles: []fileconfiguration.Profile{
			{
				Name:        settings.ProfileName,
				Environment: settings.Environment,
				Credentials: shared.Credentials{Token: settings.Token},
			},
		},
		Environments: []fileconfiguration.Environment{
			{
				Name:     settings.Environment,
				Products: envProducts,
			},
		},
	}

	return fc, nil
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

// toEndpoint converts a serverRaw into a fileconfiguration.Endpoint
func toEndpoint(s serverRaw) fileconfiguration.Endpoint {
	ep := fileconfiguration.Endpoint{SkipTLSVerify: false}

	// parse URL; relative means default API
	u, err := url.Parse(s.URL)
	if err != nil {
		ep.Name = s.URL
		return ep
	}

	if !u.IsAbs() {
		ep.Name = "https://api.ionos.com" + s.URL
		return ep
	}

	// try extracting region: host like "nfs.de-fra.ionos.com"
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) > 2 {
		loc := strings.Join(parts[1:len(parts)-2], "/")
		ep.Location = strings.ReplaceAll(loc, "-", "/")
	}

	return ep
}

func filterPages(pages []indexPage, opts Filters) []indexPage {
	latest := make(map[string]indexPage)
	for _, p := range pages {
		if opts.CustomNames != nil {
			if custom, ok := opts.CustomNames[p.Name]; ok {
				p.Name = custom
			}
		}

		if opts.Visibility != nil && *opts.Visibility != "" && p.Visibility != *opts.Visibility {
			continue
		}
		if opts.Gate != nil && *opts.Gate != "" && p.Gate != *opts.Gate {
			continue
		}
		if opts.Whitelist != nil && !opts.Whitelist[p.Name] {
			continue
		}
		if opts.Blacklist != nil && opts.Blacklist[p.Name] {
			continue
		}
		if opts.Version != nil && p.Version != *opts.Version {
			continue
		}

		prev, exists := latest[p.Name]
		if !exists || compareVersions(prev.Version, p.Version) {
			latest[p.Name] = p
		}
	}

	result := make([]indexPage, 0, len(latest))
	for _, v := range latest {
		result = append(result, v)
	}
	return result
}

// compareVersions returns true if v1 < v2
func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		n1, err1 := strconv.Atoi(parts1[i])
		n2, err2 := strconv.Atoi(parts2[i])
		if err1 != nil || err2 != nil {
			return parts1[i] < parts2[i]
		}
		if n1 != n2 {
			return n1 < n2
		}
	}
	return len(parts1) < len(parts2)
}
