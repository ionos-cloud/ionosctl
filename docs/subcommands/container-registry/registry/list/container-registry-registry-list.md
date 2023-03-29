---
description: List all Registries
---

# ContainerRegistryRegistryList

## Usage

```text
ionosctl container-registry registry list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List all managed container registries for your account

## Options

```text
      --cols strings   Set of columns to be printed on output 
                       Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -n, --name string    Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive
```

## Examples

```text
ionosctl container-registry registry list
```

