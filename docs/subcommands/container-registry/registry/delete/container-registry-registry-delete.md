---
description: Delete a Registry
---

# ContainerRegistryRegistryDelete

## Usage

```text
ionosctl container-registry registry delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Delete a Registry.

## Options

```text
  -a, --all                  Response delete all registries
      --cols strings         Set of columns to be printed on output 
                             Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -i, --registry-id string   Specify the Registry ID (required)
```

## Examples

```text
ionosctl container-registry registry delete --id [REGISTRY_ID]
```

