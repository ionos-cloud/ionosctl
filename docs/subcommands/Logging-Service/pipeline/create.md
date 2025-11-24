---
description: "Create a logging pipeline"
---

# LoggingServicePipelineCreate

## Usage

```text
ionosctl logging-service pipeline create [flags]
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

Create a logging pipeline

## Options

```text
  -u, --api-url string              Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
      --json-properties string      Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example     If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --limit int                   Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string             Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx (default "de/txl")
      --log-labels strings          Sets the labels for the pipeline log
      --log-protocol string         Sets the protocol for the pipeline log. Can be one of: http, tcp
      --log-retention-time string   Sets the retention time in days for the pipeline log. Can be one of: 7, 14, 30 (default "30")
      --log-source string           Sets the source for the pipeline log. Can be one of: docker, systemd, generic, kubernetes
      --log-tag string              Sets the tag for the pipeline log
      --log-type string             Sets the destination type for the pipeline log (default "loki")
  -n, --name string                 Sets the name of the pipeline
      --no-headers                  Don't print table headers when table output is used
      --offset int                  Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string               Desired output format [text|json|api-json] (default "text")
      --query string                JMESPath query string to filter the output
  -q, --quiet                       Quiet output
  -v, --verbose count               Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl logging-service pipeline create --json-properties PATH_TO_FILE
ionosctl logging-service pipeline create --json-properties-example
ionosctl logging-service pipeline create --name NAME --log-tag LOG_TAG --log-source LOG_SOURCE --log-protocol
LOG_PROTOCOL --log-retention-time LOG_RETENTION_TIMES
```

