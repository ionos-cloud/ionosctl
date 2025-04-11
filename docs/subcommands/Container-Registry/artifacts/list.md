---
description: "List registry or repository artifacts"
---

# ContainerRegistryArtifactsList

## Usage

```text
ionosctl container-registry artifacts list [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `artifacts` command:

```text
[a art artifact]
```

For `list` command:

```text
[l ls]
```

## Description

List all artifacts in a registry or repository

## Options

```text
  -a, --all                  List all artifacts in the registry
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Repository PushCount PullCount LastPushed TotalVulnerabilities FixableVulnerabilities MediaType URN RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings      Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -M, --max-results int32    Maximum number of results to display (default 100)
      --no-headers           Don't print table headers when table output is used
      --order-by string      Limits results to those containing a matching value for a specific property. Can be one of: -pullcount, -pushcount, -lastPush, -lastPull, -lastScan, -vulnTotalCount, -vulnFixableCount, pullCount, pushCount, lastPush, lastPull, lastScan, vulnTotalCount, vulnFixableCount (default "-pullcount")
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
      --repository string    Name of the repository to list artifacts from
  -t, --timeout duration     Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl container-registry artifacts list
```

