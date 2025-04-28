---
description: "Generate sample endpoints YAML config"
---

# Cfggen

## Usage

```text
ionosctl cfggen [flags]
```

## Description

Generate a YAML file aggregating all product endpoint information
from the public OpenAPI index. This command prints the config to stdout.

You can filter by version or specific API names.


## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --blacklist strings   Comma-separated list of API names to exclude
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
      --version string      Filter by spec version (e.g. v1)
      --whitelist strings   Comma-separated list of API names to include
```

## Examples

```text

# Generate all v1 public GA endpoints
ionosctl endpoints generate --version=v1

# Include only vpn and psql APIs, exclude billing
ionosctl endpoints generate --version=v1 \
  --whitelist=vpn,psql --blacklist=billing

```

