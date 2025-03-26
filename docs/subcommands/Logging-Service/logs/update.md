---
description: "Update a log from a logging pipeline"
---

# LoggingServiceLogsUpdate

## Usage

```text
ionosctl logging-service logs update [flags]
```

## Description

Update a log from a logging pipeline

## Options

```text
  -u, --api-url string              Override default host URL (default "https://logging.de-txl.ionos.com")
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Tag Source Protocol Public Destinations] (default [Tag,Source,Protocol,Public,Destinations])
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
  -l, --location string             Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit
      --log-labels strings          Sets the labels for the pipeline log
      --log-protocol string         Sets the protocol for the pipeline log. Can be one of: http, tcp
      --log-retention-time string   Sets the retention time in days for the pipeline log. Can be one of: 7, 14, 30 (default "30")
      --log-source string           Sets the source for the pipeline log. Can be one of: docker, systemd, generic, kubernetes
      --log-tag string              The tag of the pipeline log that you want to update (required)
      --log-type string             Sets the destination type for the pipeline log (default "loki")
      --new-log-tag string          The new tag for the pipeline log
      --no-headers                  Don't print table headers when table output is used
  -o, --output string               Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string          The ID of the logging pipeline (required)
  -q, --quiet                       Quiet output
  -t, --timeout int                 Timeout in seconds for polling the request (default 60)
  -v, --verbose                     Print step-by-step process when running command
  -w, --wait                        Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl logging-service logs update --pipeline-id ID --log-tag TAG --log-source SOURCE --log
-protocol PROTOCOL
```

