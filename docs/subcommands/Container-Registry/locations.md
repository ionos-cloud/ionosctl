---
description: "List all Registries Locations"
---

# ContainerRegistryLocations

## Usage

```text
ionosctl container-registry locations [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `locations` command:

```text
[location loc l locs]
```

## Description

List all managed container registries locations for your account

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [LocationId]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry locations
```

