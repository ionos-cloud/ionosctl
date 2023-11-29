---
description: "Create a registry"
---

# ContainerRegistryRegistryCreate

## Usage

```text
ionosctl container-registry registry create [flags]
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

For `create` command:

```text
[c]
```

## Description

Create a registry to hold container images or OCI compliant artifacts

## Options

```text
  -u, --api-url string                             Override default host url (default "https://api.ionos.com")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day using RFC3339 format
  -h, --help                                       Print usage
      --location string                            Specify the location of the registry (required)
  -n, --name string                                Specify the name of the registry (required)
      --no-headers                                 Don't print table headers when table output is used
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                                      Quiet output
  -v, --verbose                                    Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry create
```

