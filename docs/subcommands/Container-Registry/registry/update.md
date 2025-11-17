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
  -u, --api-url string                             Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                                  Level of detail for response objects (default 1)
      --filters strings                            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                      Force command to execute without user input
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day
  -h, --help                                       Print usage
      --limit int                                  Maximum number of items to return per request (default 50)
      --no-headers                                 Don't print table headers when table output is used
      --offset int                                 Number of items to skip before starting to collect the results
      --order-by string                            Property to order the results by
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
      --query string                               JMESPath query string to filter the output
  -q, --quiet                                      Quiet output
  -i, --registry-id string                         Specify the Registry ID (required)
  -v, --verbose count                              Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry registry update --id [REGISTRY_ID] --garbage-collection-schedule-days Monday
```

