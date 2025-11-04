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
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name LastSeverity ArtifactCount PullCount PushCount LastPushedAt LastPulledAt URN] (default [Id,Name,LastSeverity,ArtifactCount,PullCount,PushCount])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -F, --filters strings      Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
      --order-by string      Limits results to those containing a matching value for a specific property. Can be one of: -lastPush, -lastPull, -artifactCount, -pullCount, -pushCount, name, lastPush, lastPull, artifactCount, pullCount, pushCount (default "-lastPush")
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry list
```

