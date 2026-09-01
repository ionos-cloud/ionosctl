---
description: "Replace a registry's mutable properties (PUT)"
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

Replace the properties of an existing registry (HTTP PUT). Unlike 'update' (PATCH), this sends a full properties object, so any garbage-collection or feature field you omit is reset to its default rather than preserved.

Note: --name and --location identify an existing registry and cannot be changed by this call (the location is immutable; renaming a registry is not supported). Passing values that differ from the current ones will be rejected by the API. To only tweak the garbage-collection schedule, prefer 'container-registry registry update'.

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
      --location string                            The registry location, e.g. de/txl. Must match the existing location - it is immutable (required)
  -n, --name string                                The registry name. Must match the existing registry's name - it cannot be renamed via this call (required)
      --no-headers                                 Don't print table headers when table output is used
      --offset int                                 Number of items to skip before starting to collect the results
      --order-by string                            Property to order the results by
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
      --query string                               JMESPath query string to filter the output
  -q, --quiet                                      Quiet output
  -i, --registry-id string                         The unique ID of the registry to replace (required)
  -t, --timeout int                                Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                              Increase verbosity level [-v, -vv, -vvv]
      --vulnerability-scanning                     Enable vulnerability scanning of pushed artifacts. This is a paid add-on (default true)
  -w, --wait                                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Replace a registry, setting a single GC day and a fixed UTC GC time
ionosctl container-registry registry replace --id REGISTRY_ID --name my-registry --location de/txl --garbage-collection-schedule-days Monday --garbage-collection-schedule-time "01:00:00Z"

# Replace and disable vulnerability scanning
ionosctl container-registry registry replace --id REGISTRY_ID --name my-registry --location de/txl --vuln-scan=false
```

