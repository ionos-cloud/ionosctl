---
description: "Replace a registry"
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
  -u, --api-url string                             Override default host url (default "https://api.ionos.com")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                                      Force command to execute without user input
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day
  -h, --help                                       Print usage
      --location string                            Specify the location of the registry (required)
  -n, --name string                                Specify the name of the registry (required)
      --no-headers                                 When using text output, don't print headers
  -o, --output string                              Desired output format [text|json] (default "text")
  -q, --quiet                                      Quiet output
  -i, --registry-id string                         Specify the Registry ID (required)
  -v, --verbose                                    Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry replace --id [REGISTRY_ID] --name [REGISTRY_NAME] --location [REGISTRY_LOCATION]
```

