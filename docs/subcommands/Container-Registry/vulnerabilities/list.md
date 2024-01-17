---
description: "Retrieve vulnerabilities"
---

# ContainerRegistryVulnerabilitiesList

## Usage

```text
ionosctl container-registry vulnerabilities list [flags]
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

Retrieve all vulnerabilities from an artifact

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --artifact-id string   ID/digest of the artifact
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id DataSource Score Severity Fixable PublishedAt UpdatedAt Affects Description Recommendations References Href]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings      Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -M, --max-results int32    Maximum number of results to display (default 100)
      --no-headers           Don't print table headers when table output is used
      --order-by string      Limits results to those containing a matching value for a specific property. Can be one of: -score, -severity, -publishedAt, -updatedAt, -fixable, score, severity, publishedAt, updatedAt, fixable (default "-score")
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
      --repository string    Name of the repository to retrieve artifact from
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry vulnerabilities list
```

