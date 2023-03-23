---
description: Update the properties of a registry
---

# ContainerRegistryRegistryUpdate

## Usage

```text
ionosctl container-registry registry update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Update the "garbageCollectionSchedule" time and days of the week for runs of a registry

## Options

```text
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day
  -i, --registry-id string                         Specify the Registry ID (required)
```

## Examples

```text
ionosctl container-registry registry update --id [REGISTRY_ID]
```

