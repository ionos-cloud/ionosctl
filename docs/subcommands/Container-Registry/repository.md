---
description: "NOTE: This command's behavior will be replaced by `ionosctl container-registry repository delete` in the future. Please use that command instead.
Delete all repository contents."
---

# ContainerRegistryRepository

## Usage

```text
ionosctl container-registry repository [flags]
```

```text
ionosctl container-registry repository [command]
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

## Description

NOTE: This command's behavior will be replaced by `ionosctl container-registry repository delete` in the future. Please use that command instead.
Delete all repository contents. The registry V2 API allows manifests and blobs to be deleted individually but it is not possible to remove an entire repository. This operation is provided for convenience

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -n, --name string          Name of the repository to delete
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry repository-delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl container-registry repository delete](delete.md) | Delete all repository contents. |
| [ionosctl container-registry repository get](get.md) | Retrieve a repository. |
| [ionosctl container-registry repository list](list.md) | Retrieve all repositories. |

