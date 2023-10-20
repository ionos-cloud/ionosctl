---
description: "Update the properties of a registry"
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
  -u, --api-url string                             Override default host url (default "https://api.ionos.com")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                                      Force command to execute without user input
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day
  -h, --help                                       Print usage
      --no-headers                                 Don't print table headers when table output is used
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                                      Quiet output
  -i, --registry-id string                         Specify the Registry ID (required)
  -v, --verbose                                    Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry update --id [REGISTRY_ID]
```

