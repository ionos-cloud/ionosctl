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

Create a new Container Registry instance to hold Docker images and OCI artifacts.

The --name becomes the globally-unique hostname prefix and must be available across all IONOS customers, so check it first with 'container-registry name --name <name>'. The --location is fixed at creation and cannot be changed later (use 'container-registry locations' to list valid IDs, e.g. de/txl).

Garbage collection (--garbage-collection-schedule-days / --garbage-collection-schedule-time) is a recurring maintenance run that reclaims storage from untagged and deleted artifacts; pick a low-traffic window. Vulnerability scanning (--vuln-scan) is a paid add-on, enabled by default.

Once the registry is AVAILABLE, authenticate with 'docker login <hostname>' using a token created via 'container-registry token create'.

## Options

```text
  -u, --api-url string                             Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                                  Level of detail for response objects (default 1)
  -F, --filters strings                            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                      Force command to execute without user input
      --garbage-collection-schedule-days strings   Weekly days on which garbage collection runs to reclaim storage from untagged/deleted artifacts. Comma-separated full weekday names (Monday...Sunday). Defaults to a random day Mon-Fri (default Random (Mon-Fri 10:00-16:00))
      --garbage-collection-schedule-time string    UTC time of day at which garbage collection runs, as an RFC3339 partial-time. e.g. "16:00:00Z" or "01:23:00+00:00". Defaults to a random hour in 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -h, --help                                       Print usage
      --limit int                                  Maximum number of items to return per request (default 50)
  -l, --location string                            The location that will host the registry, e.g. de/txl. Fixed at creation - it cannot be changed later. See 'container-registry locations' for valid IDs (required)
  -n, --name string                                The name of the registry. Becomes the hostname prefix and must be globally unique across all IONOS customers. Lowercase letters, digits and dashes only, 3-63 chars, starting with a letter (regex ^[a-z][-a-z0-9]{1,61}[a-z0-9]$). Check availability with 'container-registry name' (required)
      --no-headers                                 Don't print table headers when table output is used
      --offset int                                 Number of items to skip before starting to collect the results
      --order-by string                            Property to order the results by
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
      --query string                               JMESPath query string to filter the output
  -q, --quiet                                      Quiet output
  -t, --timeout int                                Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                              Increase verbosity level [-v, -vv, -vvv]
      --vulnerability-scanning                     Enable vulnerability scanning of pushed artifacts. This is a paid add-on; enabled by default (default true)
  -w, --wait                                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a registry with defaults (vulnerability scanning on, random GC window Mon-Fri 10:00-16:00)
ionosctl container-registry registry create --name my-registry --location de/txl

# Create a registry with an explicit weekend GC window and vulnerability scanning disabled
ionosctl container-registry registry create --name my-registry --location de/txl --garbage-collection-schedule-days Saturday,Sunday --garbage-collection-schedule-time "02:00:00Z" --vuln-scan=false
```

