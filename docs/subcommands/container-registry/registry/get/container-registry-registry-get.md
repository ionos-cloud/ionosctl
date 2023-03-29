---
description: Get Properties of a Registry
---

# ContainerRegistryRegistryGet

## Usage

```text
ionosctl container-registry registry get [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `registry` command:

```text
[reg registries r]
```

For `get` command:

```text
[g]
```

## Description

Get Properties of a single Registry

## Options

```text
      --cols strings         Set of columns to be printed on output 
                             Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -i, --registry-id string   Registry ID (required)
```

## Examples

```text
ionosctl container-registry registry get --id [REGISTRY_ID]
```

