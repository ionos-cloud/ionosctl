---
description: "Update a registry's garbage-collection schedule (PATCH)"
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

Update the garbage-collection schedule of an existing registry (HTTP PATCH). Garbage collection is the recurring maintenance run that reclaims storage occupied by untagged and deleted artifacts.

Only the days and time of the schedule can be changed here; the registry name, location and features are not touched. Set --garbage-collection-schedule-days to the weekly run days and --garbage-collection-schedule-time to the UTC time of day.

## Options

```text
  -u, --api-url string                             Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                                  Level of detail for response objects (default 1)
  -F, --filters strings                            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                      Force command to execute without user input
      --garbage-collection-schedule-days strings   Weekly days on which garbage collection runs. Comma-separated full weekday names (Monday...Sunday)
      --garbage-collection-schedule-time string    UTC time of day for garbage collection, as an RFC3339 partial-time, e.g. "01:23:00+00:00"
  -h, --help                                       Print usage
      --limit int                                  Maximum number of items to return per request (default 50)
      --no-headers                                 Don't print table headers when table output is used
      --offset int                                 Number of items to skip before starting to collect the results
      --order-by string                            Property to order the results by
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
      --query string                               JMESPath query string to filter the output
  -q, --quiet                                      Quiet output
  -i, --registry-id string                         The unique ID of the registry to update (required)
  -t, --timeout int                                Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                              Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Move garbage collection to Monday
ionosctl container-registry registry update --id REGISTRY_ID --garbage-collection-schedule-days Monday

# Run garbage collection on the weekend, early morning UTC
ionosctl container-registry registry update --id REGISTRY_ID --garbage-collection-schedule-days Saturday,Sunday --garbage-collection-schedule-time "03:00:00Z"
```

