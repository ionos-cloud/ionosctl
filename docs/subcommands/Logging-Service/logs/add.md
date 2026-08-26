---
description: "Add a log stream to a logging pipeline"
---

# LoggingServiceLogsAdd

## Usage

```text
ionosctl logging-service logs add [flags]
```

## Description

Add a new log stream to an existing pipeline. The log is appended to the pipeline's current logs (nothing existing is replaced), so the tag must be unique within the pipeline.

Valid values:
  --log-source:   docker, systemd, kubernetes, generic
  --log-protocol: http or tcp
  --log-type:     destination backend, currently loki
  --log-retention-time: 7, 14 or 30 days

## Options

```text
  -u, --api-url string              Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Tag Source Protocol Public Destinations Labels PipelineId]
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                   Level of detail for response objects (default 1)
  -F, --filters strings             Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
      --limit int                   Maximum number of items to return per request (default 50)
  -l, --location string             Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --log-labels strings          Optional labels attached to every entry of this stream for filtering in Grafana. Comma-separated, e.g. --log-labels env=prod,team=core
      --log-protocol string         Transport used to push logs: 'http' targets the pipeline's HTTP address, 'tcp' its TCP address. Can be one of: http, tcp (required)
      --log-retention-time string   How many days logs are kept before deletion. One of: 7, 14, 30. Can be one of: 7, 14, 30 (default "30")
      --log-source string           The kind of workload producing the logs, which selects how they are parsed. One of: docker, systemd, kubernetes, generic. Can be one of: docker, systemd, generic, kubernetes (required)
      --log-tag string              Tag identifying this log stream within the pipeline; must be unique and is how you reference the log later (required)
      --log-type string             Destination backend the logs are stored in and queried from. Currently 'loki' (default "loki")
      --no-headers                  Don't print table headers when table output is used
      --offset int                  Number of items to skip before starting to collect the results
      --order-by string             Property to order the results by
  -o, --output string               Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string          The ID of the logging pipeline to add the log stream to (required)
      --query string                JMESPath query string to filter the output
  -q, --quiet                       Quiet output
  -t, --timeout int                 Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count               Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                        Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Add a systemd log stream shipped over HTTP
ionosctl logging-service logs add --location de/txl --pipeline-id ID --log-tag sys --log-source systemd --log-protocol http

# Advanced: 14-day retention with labels
ionosctl logging-service logs add --location de/txl --pipeline-id ID --log-tag k8s --log-source kubernetes --log-protocol tcp --log-retention-time 14 --log-labels env=staging
```

