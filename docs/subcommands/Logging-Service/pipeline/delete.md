---
description: "Delete a logging pipeline using its ID"
---

# LoggingServicePipelineDelete

## Usage

```text
ionosctl logging-service pipeline delete [flags]
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

## Description

Delete a logging pipeline using its ID

## Options

```text
  -a, --all                  Use this flag to delete all logging pipelines
  -u, --api-url string       Override default host URL (default "https://logging.de-txl.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -l, --location string      Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string   The ID of the logging pipeline you want to delete (required)
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service pipeline delete --pipeline-id ID
```

