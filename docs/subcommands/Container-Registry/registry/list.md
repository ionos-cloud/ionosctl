---
description: "List all Registries"
---

# ContainerRegistryRegistryList

## Usage

```text
ionosctl container-registry registry list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List all managed container registries for your account

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --name string      Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry list
```

