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
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id DataSource Score Severity Fixable PublishedAt UpdatedAt Affects Description Recommendations References Href]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -v, --verbose                   Print step-by-step process when running command
      --vulnerability-id string   Vulnerability ID
```

## Examples

```text
ionosctl container-registry vulnerabilities get
```

