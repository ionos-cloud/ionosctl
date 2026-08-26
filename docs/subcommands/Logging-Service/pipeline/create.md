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

Create a logging pipeline together with its first log stream. On success the pipeline's ingestion addresses (TCP/HTTP) and Grafana address are returned; generate a key with 'pipeline key' before shipping logs.

There are two ways to define the pipeline:
  1. Flags: --name plus a single log described by --log-tag, --log-source, --log-protocol (and optionally --log-type, --log-retention-time, --log-labels). This creates a pipeline with exactly one log.
  2. --json-properties: a path to a JSON file describing the full pipeline, which lets you define MULTIPLE logs at once. Use --json-properties-example to print a ready-to-edit template.

A pipeline must contain at least one log, which is why source, tag and protocol are required in flag mode.

Valid values:
  --log-source:   docker, systemd, kubernetes, generic (the kind of workload the log comes from)
  --log-protocol: http or tcp (how your shipper connects to the ingestion endpoint)
  --log-type:     the destination backend; currently loki
  --log-retention-time: 7, 14 or 30 days

Note: --name and --log-source are normalised to lower-case before being sent.

## Options

```text
  -u, --api-url string              Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Id Name GrafanaAddress TCPAddress HTTPAddress CreatedDate State]
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                   Level of detail for response objects (default 1)
  -F, --filters strings             Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
      --json-properties string      Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example     If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --limit int                   Maximum number of items to return per request (default 50)
  -l, --location string             Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --log-labels strings          Optional labels attached to every entry of this log stream, useful for filtering in Grafana. Comma-separated, e.g. --log-labels env=prod,team=core
      --log-protocol string         Transport your shipper uses to push logs to the ingestion endpoint: 'http' targets the pipeline's HTTP address, 'tcp' its TCP address. Can be one of: http, tcp
      --log-retention-time string   How many days logs are kept before being deleted. One of: 7, 14, 30. Can be one of: 7, 14, 30 (default "30")
      --log-source string           The kind of workload producing the logs, which selects how they are parsed. One of: docker, systemd, kubernetes, generic (use 'generic' for anything that does not fit the others). Can be one of: docker, systemd, generic, kubernetes
      --log-tag string              Tag that identifies this log stream within the pipeline; used later to reference the log (e.g. in 'logs get/update/remove') and to route entries
      --log-type string             Destination backend the logs are stored in and queried from. Currently 'loki' (default "loki")
  -n, --name string                 Human-readable name of the pipeline, shown in listings and Grafana. Normalised to lower-case
      --no-headers                  Don't print table headers when table output is used
      --offset int                  Number of items to skip before starting to collect the results
      --order-by string             Property to order the results by
  -o, --output string               Desired output format [text|json|api-json] (default "text")
      --query string                JMESPath query string to filter the output
  -q, --quiet                       Quiet output
  -t, --timeout int                 Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count               Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                        Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a pipeline with a single Kubernetes log stream shipped over TCP
ionosctl logging-service pipeline create --location de/txl --name my-pipeline --log-tag k8s --log-source kubernetes --log-protocol tcp

# Advanced: pin retention and attach labels, shipping Docker logs over HTTP
ionosctl logging-service pipeline create --location de/txl --name app-logs --log-tag docker-app --log-source docker --log-protocol http --log-retention-time 30 --log-labels env=prod,team=core

# Create from a JSON file (allows multiple log streams in one pipeline)
ionosctl logging-service pipeline create --location de/txl --json-properties ./pipeline.json

# Print a JSON template you can edit and feed back via --json-properties
ionosctl logging-service pipeline create --json-properties-example
```

