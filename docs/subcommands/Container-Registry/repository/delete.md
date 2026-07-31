---
description: "Delete all contents of a repository"
---

# ContainerRegistryRepositoryDelete

## Usage

```text
ionosctl container-registry repository delete [flags]
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

For `delete` command:

```text
[d del]
```

## Description

Delete all contents (every artifact and tag) of a repository, effectively removing the repository from the registry.

The registry V2 API only allows manifests and blobs to be deleted individually, not whole repositories, so this command performs that cleanup for you as a convenience. This is irreversible.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name LastSeverity ArtifactCount PullCount PushCount LastPushedAt LastPulledAt URN]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -n, --name string          Name of the repository to delete (the path in <hostname>/<repository>)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   The unique ID of the registry the repository belongs to
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl container-registry repository delete --registry-id REGISTRY_ID --name REPOSITORY_NAME
```

