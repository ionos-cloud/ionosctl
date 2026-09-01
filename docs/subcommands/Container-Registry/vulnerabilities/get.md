---
description: "Get a single vulnerability finding by ID"
---

# ContainerRegistryVulnerabilitiesGet

## Usage

```text
ionosctl container-registry vulnerabilities get [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `vulnerabilities` command:

```text
[v vuln vulnerability]
```

## Description

Get the full details of one vulnerability finding by its ID (e.g. a CVE identifier), including CVSS score, severity, affected packages/versions, remediation recommendations and references. Find IDs via 'container-registry vulnerabilities list'.

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id DataSource Score Severity Fixable PublishedAt UpdatedAt Affects Description Recommendations References Href]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --limit int                 Maximum number of items to return per request (default 50)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --vulnerability-id string   The ID of the vulnerability finding to retrieve (e.g. a CVE identifier, as shown by 'vulnerabilities list')
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl container-registry vulnerabilities get --vulnerability-id VULNERABILITY_ID
```

