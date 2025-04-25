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
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id DataSource Score Severity Fixable PublishedAt UpdatedAt Affects Description Recommendations References Href]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -t, --timeout duration          Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                   Print step-by-step process when running command
      --vulnerability-id string   Vulnerability ID
  -w, --wait                      Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl container-registry vulnerabilities get
```

