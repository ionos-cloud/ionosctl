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
      --cols strings    Set of columns to be printed on output 
                        Available columns: [LocationId]
  -c, --config string   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force           Force command to execute without user input
  -h, --help            Print usage
      --no-headers      Don't print table headers when table output is used
  -o, --output string   Desired output format [text|json|api-json] (default "text")
  -q, --quiet           Quiet output
  -v, --verbose         Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry locations
```

