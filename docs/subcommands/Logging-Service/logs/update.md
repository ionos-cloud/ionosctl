---
description: "Update a log stream in a logging pipeline"
---

# LoggingServiceLogsUpdate

## Usage

```text
ionosctl logging-service logs update [flags]
```

## Description

Update one log stream, selected by its --log-tag, and patch it back into the pipeline. Only the flags you pass are changed; every other attribute of the log is preserved. You can rename the tag (--new-log-tag), change the --log-source or --log-protocol, and adjust the destination backend (--log-type) and retention (--log-retention-time).

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
      --log-labels strings          Labels for the log stream. Comma-separated, e.g. --log-labels env=prod,team=core
      --log-protocol string         New protocol for the log stream. One of: http, tcp. Leave unset to keep the current protocol. Can be one of: http, tcp
      --log-retention-time string   How many days logs are kept before deletion. One of: 7, 14, 30. Can be one of: 7, 14, 30 (default "30")
      --log-source string           New source for the log stream. One of: docker, systemd, kubernetes, generic. Leave unset to keep the current source. Can be one of: docker, systemd, generic, kubernetes
      --log-tag string              Tag of the log stream to update (identifies which log within the pipeline) (required)
      --log-type string             Destination backend the logs are stored in and queried from. Currently 'loki' (default "loki")
      --new-log-tag string          Rename the log stream's tag. Leave unset to keep the current tag
      --no-headers                  Don't print table headers when table output is used
      --offset int                  Number of items to skip before starting to collect the results
      --order-by string             Property to order the results by
  -o, --output string               Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string          The ID of the logging pipeline containing the log stream (required)
      --query string                JMESPath query string to filter the output
  -q, --quiet                       Quiet output
  -t, --timeout int                 Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count               Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                        Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Extend a log stream's retention to 30 days
ionosctl logging-service logs update --location de/txl --pipeline-id ID --log-tag k8s --log-retention-time 30
```

