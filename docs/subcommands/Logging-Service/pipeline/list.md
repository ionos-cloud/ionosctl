---
description: "Retrieve logging pipelines"
---

# LoggingServicePipelineList

## Usage

```text
ionosctl logging-service pipeline list [flags]
```

## Aliases

For `logging-service` command:

```text
[log-svc]
```

For `pipeline` command:

```text
[p pipelines]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve logging pipelines

## Options

```text
  -u, --api-url string    Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://logging.de-txl.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service pipeline list
```

