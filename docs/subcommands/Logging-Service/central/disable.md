---
description: "Disable CentralLogging"
---

# LoggingServiceCentralDisable

## Usage

```text
ionosctl logging-service central disable [flags]
```

## Aliases

For `logging-service` command:

```text
[log-svc]
```

For `central` command:

```text
[c]
```

For `disable` command:

```text
[d]
```

## Description

Disable CentralLogging

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx (default "de/txl")
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service central disable --location de/txl
```

