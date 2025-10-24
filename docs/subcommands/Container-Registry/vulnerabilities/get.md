---
description: "Retrieve a vulnerability"
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

Retrieve a vulnerability

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id DataSource Score Severity Fixable PublishedAt UpdatedAt Affects Description Recommendations References Href]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --vulnerability-id string   Vulnerability ID
```

## Examples

```text
ionosctl container-registry vulnerabilities get
```

