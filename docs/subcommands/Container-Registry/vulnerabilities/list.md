---
description: "List vulnerabilities found in an artifact"
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

List all vulnerability findings for a single artifact, identified by registry, repository and artifact digest. Requires the registry's vulnerabilityScanning feature to be enabled; if it is off, no findings are returned. Each row shows the CVE, score, severity and whether a fix exists.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --artifact-id string   Content digest of the artifact to scan for findings, e.g. sha256:12ab...
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id DataSource Score Severity Fixable PublishedAt UpdatedAt Affects Description Recommendations References Href]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   The unique ID of the registry the artifact belongs to
      --repository string    Name of the repository that holds the artifact
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl container-registry vulnerabilities list --registry-id REGISTRY_ID --repository my-app --artifact-id sha256:DIGEST
```

