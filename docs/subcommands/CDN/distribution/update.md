---
description: "Update a distribution's domain, certificate binding, or routing rules"
---

# CdnDistributionUpdate

## Usage

```text
ionosctl cdn distribution update [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `update` command:

```text
[u]
```

## Description

Update an existing CDN distribution. Only the properties you pass are changed: the command first GETs the current distribution, overlays the flags you set (--domain, --certificate-id, --routing-rules), and PUTs the result back, so unspecified properties are preserved (a PATCH-like behavior).

Note that --routing-rules REPLACES the entire rule list; there is no way to edit a single rule in place. To modify one rule, fetch the current rules with 'ionosctl cdn ds rr get --distribution-id <id> -o json', edit the JSON, and pass the full array back. Provide 1-25 rules. See 'ionosctl cdn ds create --routing-rules-example' for the JSON format and field meanings (scheme, upstream host/caching/waf/rateLimitClass/sniMode/geoRestrictions).

## Options

```text
  -u, --api-url string           Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cdn' and env var 'IONOS_API_URL' (default "https://cdn.%s.ionos.com")
      --certificate-id string    Certificate Manager UUID used to terminate HTTPS for the domain. Omit for an HTTP-only distribution
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                Level of detail for response objects (default 1)
  -i, --distribution-id string   The ID of the distribution you want to update (required)
      --domain string            The public hostname this distribution serves, e.g. cdn.example.com. Must be a valid, unique domain (2-253 chars)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Maximum number of items to return per request (default 50)
  -l, --location string          Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --routing-rules string     Routing rules as a JSON array (inline string or a path to a .json file). Each rule maps a path prefix + scheme to an upstream origin (host, caching, waf, rateLimitClass, sniMode, geoRestrictions). 1-25 rules. See --routing-rules-example for the format
      --routing-rules-example    Print a ready-to-edit routing-rules JSON template and exit (does not create anything). Redirect to a file, edit it, then pass it via --routing-rules
  -t, --timeout int              Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                     Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rebind the distribution to a new HTTPS certificate
ionosctl cdn ds update --distribution-id <id> --certificate-id 5a029f4a-72e5-11ec-90d6-0242ac120003

# Replace all routing rules from a file
ionosctl cdn ds update --distribution-id <id> --routing-rules rules.json
```

