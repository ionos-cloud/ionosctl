---
description: "enable a CentralMonitoring"
---

# MonitoringCentralEnable

## Usage

```text
ionosctl monitoring central enable [flags]
```

## Aliases

For `central` command:

```text
[c]
```

For `enable` command:

```text
[e]
```

## Description

enable a CentralMonitoring

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'monitoring' and env var 'IONOS_API_URL' (default "https://monitoring.%s.ionos.com")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/bhx, gb/lhr, fr/par, us/mci (default "de/fra")
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl monitoring central enable --location de/txl
```

