---
description: "Get Properties of a Registry"
---

# ContainerRegistryRegistryGet

## Usage

```text
ionosctl container-registry registry get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get Properties of a single Registry

## Options

```text
      --cols strings         Set of columns to be printed on output 
                             Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --registry-id string   Registry ID (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry get --id [REGISTRY_ID]
```

