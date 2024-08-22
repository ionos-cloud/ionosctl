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
  -u, --api-url string   Override default host url (default "logging.de-txl.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service pipeline list
```

