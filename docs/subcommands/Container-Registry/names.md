---
description: "Check if a Registry Name is available"
---

# ContainerRegistryNames

## Usage

```text
ionosctl container-registry names [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `names` command:

```text
[check name n]
```

## Description

Check if a Registry Name is available

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --name string      Name to check availability for (required)
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry name
```

