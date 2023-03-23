---
description: Replace a registry
---

# ContainerRegistryRegistryReplace

## Usage

```text
ionosctl container-registry registry replace [flags]
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

For `replace` command:

```text
[r rep]
```

## Description

Create/replace a registry to hold container images or OCI compliant artifacts

## Options

```text
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day
      --location string                            Specify the location of the registry (required)
  -n, --name string                                Specify the name of the registry (required)
  -i, --registry-id string                         Specify the Registry ID (required)
```

## Examples

```text
ionosctl container-registry registry replace --id [REGISTRY_ID] --name [REGISTRY_NAME] --location [REGISTRY_LOCATION]
```

