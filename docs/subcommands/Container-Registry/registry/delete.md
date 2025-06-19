---
description: "Delete a Registry"
---

# ContainerRegistryRegistryDelete

## Usage

```text
ionosctl container-registry registry delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Delete a Registry.

## Options

```text
  -a, --all                  Response delete all registries
      --cols strings         Set of columns to be printed on output 
                             Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --registry-id string   Specify the Registry ID (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry registry delete --id [REGISTRY_ID]
```

