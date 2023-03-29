---
description: Delete all repository contents.
---

# ContainerRegistryRepository

## Usage

```text
ionosctl container-registry repository [flags]
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

Delete all repository contents. The registry V2 API allows manifests and blobs to be deleted individually but it is not possible to remove an entire repository. This operation is provided for convenience

## Options

```text
  -n, --name string          Name of the repository to delete
  -r, --registry-id string   Registry ID
```

## Examples

```text
ionosctl container-registry repository-delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]
```

