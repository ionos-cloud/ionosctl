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

Delete all repository contents. The registry V2 API allows manifests and blobs to be deleted individually but it is not possible to remove an entire repository. This operation is provided for convenience

## Options

```text
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
ionosctl container-registry delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]
```

