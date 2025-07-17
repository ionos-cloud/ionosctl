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
  -u, --api-url string                             Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings                               Set of columns to be printed on output 
                                                   Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                                      Force command to execute without user input
      --garbage-collection-schedule-days strings   Specify the garbage collection schedule days. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default Random (Mon-Fri 10:00-16:00))
      --garbage-collection-schedule-time string    Specify the garbage collection schedule time of day using RFC3339 format. i.e. "16:00:00Z". Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -h, --help                                       Print usage
      --location string                            Specify the location of the registry (required)
  -n, --name string                                Specify the name of the registry (required)
      --no-headers                                 Don't print table headers when table output is used
  -o, --output string                              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                                      Quiet output
  -v, --verbose                                    Print step-by-step process when running command
      --vulnerability-scanning                     Enable/disable vulnerability scanning (this is a paid add-on) (default true)
```

## Examples

```text
ionosctl container-registry registry create
```

