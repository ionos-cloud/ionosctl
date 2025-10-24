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
  -u, --api-url string   Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --name string      Name to check availability for (required)
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry name
```

