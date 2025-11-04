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
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx (default "de/txl")
      --no-headers        Don't print table headers when table output is used
      --offset int        Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl logging-service pipeline list
```

