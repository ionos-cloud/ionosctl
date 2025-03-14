---
description: "Retrieve all repositories."
---

# ContainerRegistryRepositoryList

## Usage

```text
ionosctl container-registry repository list [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `repository` command:

```text
[rd del repo rep-del repository-delete]
```

For `list` command:

```text
[ls l]
```

## Description

Retrieve all repositories in a registry.

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name LastSeverity ArtifactCount PullCount PushCount LastPushedAt LastPulledAt URN] (default [Id,Name,LastSeverity,ArtifactCount,PullCount,PushCount])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings      Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -M, --max-results int32    Maximum number of results to display (default 100)
      --no-headers           Don't print table headers when table output is used
      --order-by string      Limits results to those containing a matching value for a specific property. Can be one of: -lastPush, -lastPull, -artifactCount, -pullCount, -pushCount, name, lastPush, lastPull, artifactCount, pullCount, pushCount (default "-lastPush")
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -t, --timeout int          Timeout in seconds for polling the request (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl container-registry list
```

