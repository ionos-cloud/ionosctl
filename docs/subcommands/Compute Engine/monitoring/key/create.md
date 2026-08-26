---
description: "Generate (rotate) a pipeline's ingest key"
---

# MonitoringKeyCreate

## Usage

```text
ionosctl monitoring key create [flags]
```

## Aliases

For `key` command:

```text
[k]
```

For `create` command:

```text
[post c generate]
```

## Description

Generate a fresh ingest key for a pipeline and print it. This is the only way to recover a key after creation, since it is shown just once.

Rotating is destructive to the old key: the previous key is invalidated immediately, so update every agent pushing to this pipeline with the new value. Agents send the key as the APIKEY header when writing to <HttpEndpoint>/api/v1/push (Prometheus remote-write).

The command prints only the raw key so it can be captured directly, e.g. piped into a secret store.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'monitoring' and env var 'IONOS_API_URL' (default "https://monitoring.%s.ionos.com")
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/bhx, gb/lhr, fr/par, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string   The ID of the pipeline whose ingest key should be rotated (from 'pipeline list')
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rotate the key and print it
ionosctl monitoring key create --location de/txl --pipeline-id PIPELINE_ID

# Capture the new key into a variable
KEY=$(ionosctl monitoring key create --location de/txl --pipeline-id PIPELINE_ID)
```

