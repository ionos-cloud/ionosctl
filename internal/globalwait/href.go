package globalwait

import (
	"encoding/json"
	neturl "net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"
)

// extractHref extracts the top-level "href" field from sourceData.
// Returns empty string if sourceData is a list (has "items" key), has no href,
// or cannot be parsed.
func extractHref(sourceData any) string {
	m := extractMap(sourceData)
	if m == nil {
		return ""
	}
	if href, ok := m["href"].(string); ok {
		return href
	}
	return ""
}

// extractID extracts the top-level "id" field from sourceData.
// Used to build resource URLs for APIs that don't include href in responses
// (e.g. postgres-v1, mongo).
func extractID(sourceData any) string {
	m := extractMap(sourceData)
	if m == nil {
		return ""
	}
	if id, ok := m["id"].(string); ok {
		return id
	}
	return ""
}

func extractMap(sourceData any) map[string]any {
	b, err := json.Marshal(sourceData)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	// Skip list/collection responses
	if _, hasItems := m["items"]; hasItems {
		return nil
	}
	return m
}

// isActionEndpoint returns true for action endpoints that don't support GET
// (e.g. server start/stop, database restore, DNS zone transfer).
func isActionEndpoint(href string) bool {
	u, err := neturl.Parse(href)
	if err != nil {
		return false
	}
	parts := strings.Split(strings.TrimRight(u.Path, "/"), "/")
	if len(parts) == 0 {
		return false
	}
	last := parts[len(parts)-1]
	switch last {
	case "start", "stop", "reboot", "suspend", "resume", "restore", "transfer":
		return true
	}
	return false
}

// resourceAndParentURLs returns the given href plus all parent resource hrefs,
// from deepest to shallowest. For example, given:
//
//	https://api.ionos.com/cloudapi/v6/datacenters/dc1/servers/srv1/volumes/vol1
//
// it returns:
//
//	[".../volumes/vol1", ".../servers/srv1", ".../datacenters/dc1"]
//
// Non-CloudAPI hrefs (no /cloudapi/ path) return just the original href.
func resourceAndParentURLs(href string) []string {
	urls := []string{href}

	// Walk up by stripping last two path segments (resource-type/resource-id)
	current := href
	for {
		parent := parentHref(current)
		if parent == "" {
			break
		}
		urls = append(urls, parent)
		current = parent
	}

	return urls
}

// parentHref strips the last two path segments (resource-type/id) to get the
// parent resource href. Returns "" if there's no valid parent.
//
// Works with both CloudAPI and regional API URL structures:
//   - CloudAPI: https://api.ionos.com/cloudapi/v6/datacenters/dc1/servers/srv1
//   - Regional: https://vpn.de-fra.ionos.com/wireguardgateways/gw1/peers/p1
//
// Stops when stripping would leave no resource pair (type+id) after the host.
func parentHref(href string) string {
	parts := strings.Split(href, "/")

	// After splitting, first 3 parts are always: "https:", "", "host"
	// Then optional API prefix segments (e.g. "cloudapi", "v6") followed by
	// resource pairs (type/id). We need at least 2 resource pairs (4 path
	// segments after host) to have a parent.
	//
	// Find where resource path starts by skipping non-UUID/non-resource segments
	// after the host. Simpler: the candidate after stripping 2 must still have
	// at least one type/id pair (2 segments) after the host portion.
	//
	// Minimum: "https:"/""/host/type/id/type/id = 7 parts
	// After strip: "https:"/""/host/type/id = 5 parts (valid parent)
	// Next strip would give: "https:"/""/host = 3 parts (just host, not valid)
	if len(parts) < 7 {
		return ""
	}

	candidate := strings.Join(parts[:len(parts)-2], "/")

	// Candidate must end with a resource ID (last segment should look like
	// an ID, not an API path component like "v6" or "cloudapi").
	// Resource IDs are typically UUIDs or alphanumeric strings.
	lastSeg := parts[len(parts)-3] // last segment of candidate
	if !looksLikeResourceID(lastSeg) {
		return ""
	}

	return candidate
}

// uuidRegex matches standard UUID format used by IONOS APIs.
var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// looksLikeResourceID returns true if the string looks like a resource ID
// (UUID or pure numeric). Uses strict UUID regex to avoid false positives
// on hyphenated resource type names like "private-cross-connects".
func looksLikeResourceID(s string) bool {
	if s == "" {
		return false
	}
	if uuidRegex.MatchString(s) {
		return true
	}
	// Pure numeric IDs
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func buildFullURL(href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return appendDepthParam(href)
	}

	baseURL := viper.GetString(constants.ArgServerUrl)
	if baseURL == "" {
		baseURL = constants.DefaultApiURL
	}

	if !strings.HasPrefix(href, "/") {
		href = "/" + href
	}

	return appendDepthParam(strings.TrimRight(baseURL, "/") + href)
}

func appendDepthParam(rawURL string) string {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	depth := viper.GetInt(constants.FlagDepth)
	if depth <= 0 {
		depth = 1
	}
	q := u.Query()
	q.Set("depth", strconv.Itoa(depth))
	u.RawQuery = q.Encode()
	return u.String()
}

func isStructuredOutput() bool {
	switch viper.GetString(constants.ArgOutput) {
	case "json", "api-json":
		return true
	default:
		return false
	}
}
