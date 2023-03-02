---
description: List all Registries
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
[l]
```

## Description

List all managed container registries for your account

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [RegistryId DisplayName Location Hostname GarbageCollectionDays GarbageCollectionTime]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --name string      Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry list
```

