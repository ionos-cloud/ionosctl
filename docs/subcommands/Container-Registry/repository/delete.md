---
description: "Delete all repository contents."
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

Delete all repository contents.
The registry V2 API allows manifests and blobs to be deleted individually, but it is not possible to remove an entire repository. This operation is provided for convenience

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string          Name of the repository to delete
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]
```

