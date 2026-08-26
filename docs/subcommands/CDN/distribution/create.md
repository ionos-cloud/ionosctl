---
description: "Create a CDN distribution for a domain, with routing rules mapping URL prefixes to upstream origins"
---

# CdnDistributionCreate

## Usage

```text
ionosctl cdn distribution create [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `create` command:

```text
[c post]
```

## Description

Create a CDN distribution. A distribution serves a single DOMAIN and needs at least one ROUTING RULE describing where to fetch content that is not already cached.

Each routing rule matches requests whose path starts with a given prefix (e.g. "/api") and a scheme (http, https, or http/https), then forwards them to an upstream origin. Per rule you control caching, the WAF, a per-IP rate-limit class, geo-restrictions, and the SNI mode used when the CDN connects to the origin over TLS. Rules are supplied as JSON via --routing-rules; run 'ionosctl cdn ds create --routing-rules-example' to print a ready-to-edit template.

Provide --certificate-id (a Certificate Manager UUID) to terminate HTTPS for the domain; omit it for HTTP-only distributions.

Constraints (enforced by the API):
  - --domain must be a valid, unique hostname (2-253 chars), e.g. cdn.example.com.
  - Each distribution needs 1-25 routing rules; each rule's prefix is 1-128 chars and must start with "/".
  - Once AVAILABLE, point the domain's DNS (usually a CNAME) at the CDN so traffic reaches the edge.

Routing-rule JSON fields (per rule):
  - prefix:   URL path prefix to match, e.g. "/" or "/api".
  - scheme:   one of "http", "https", "http/https" (accept both).
  - upstream.host:           origin hostname to fetch uncached content from.
  - upstream.caching:        true/false; cache origin responses at the edge.
  - upstream.waf:            true/false; enable the Web Application Firewall.
  - upstream.rateLimitClass: per-IP request rate limit, one of R1, R5, R10, R25, R50, R100, R250, R500 (the number is the allowed requests/second per client IP; R1 is strictest, R500 most permissive).
  - upstream.sniMode:        "distribution" (origin cert must match the distribution's domain) or "origin" (origin cert must match upstream.host).
  - upstream.geoRestrictions: optionally EITHER {"allowList":[...]} (only these countries may access) OR {"blockList":[...]} (these countries are denied), using ISO 3166-1 alpha-2 codes (e.g. "DE", "US"). Use one list per rule, not both.

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cdn' and env var 'IONOS_API_URL' (default "https://cdn.%s.ionos.com")
      --certificate-id string   Certificate Manager UUID used to terminate HTTPS for the domain. Omit for an HTTP-only distribution
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Domain CertificateId State]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
      --domain string           The public hostname this distribution serves, e.g. cdn.example.com. Must be a valid, unique domain (2-253 chars)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
      --routing-rules string    Routing rules as a JSON array (inline string or a path to a .json file). Each rule maps a path prefix + scheme to an upstream origin (host, caching, waf, rateLimitClass, sniMode, geoRestrictions). 1-25 rules. See --routing-rules-example for the format
      --routing-rules-example   Print a ready-to-edit routing-rules JSON template and exit (does not create anything). Redirect to a file, edit it, then pass it via --routing-rules
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create an HTTP-only distribution, passing routing rules inline as JSON
ionosctl cdn ds create --domain cdn.example.com --routing-rules '[{"prefix":"/","scheme":"http/https","upstream":{"host":"origin.example.com","caching":true,"waf":false,"rateLimitClass":"R500","sniMode":"origin"}}]'

# Print an editable routing-rules template, then create an HTTPS distribution from a rules file
ionosctl cdn ds create --routing-rules-example > rules.json
ionosctl cdn ds create --domain cdn.example.com --certificate-id 5a029f4a-72e5-11ec-90d6-0242ac120003 --routing-rules rules.json
```

