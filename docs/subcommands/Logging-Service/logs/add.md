---
description: "Add a log to a logging pipeline"
---

# LoggingServiceLogsAdd

## Usage

```text
ionosctl logging-service logs add [flags]
```

## Description

Add a log to a logging pipeline

## Options

```text
  -u, --api-url string              Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Tag Source Protocol Public Destinations] (default [Tag,Source,Protocol,Public,Destinations])
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
  -l, --location string             Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit (default "de/txl")
      --log-labels strings          Sets the labels for the pipeline log
      --log-protocol string         Sets the protocol for the pipeline log. Can be one of: http, tcp (required)
      --log-retention-time string   Sets the retention time in days for the pipeline log. Can be one of: 7, 14, 30 (default "30")
      --log-source string           Sets the source for the pipeline log. Can be one of: docker, systemd, generic, kubernetes (required)
      --log-tag string              Sets the tag for the pipeline log (required)
      --log-type string             Sets the destination type for the pipeline log (default "loki")
      --no-headers                  Don't print table headers when table output is used
  -o, --output string               Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string          The ID of the logging pipeline (required)
  -q, --quiet                       Quiet output
  -v, --verbose                     Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service logs add --pipeline-id ID --log-tag TAG --log-source SOURCE --log
-protocol PROTOCOL
```

