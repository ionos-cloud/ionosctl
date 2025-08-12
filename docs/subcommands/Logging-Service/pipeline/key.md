---
description: "Generate a new key for a logging pipeline, invalidating the old one. The key is used for authentication when sending logs."
---

# LoggingServicePipelineKey

## Usage

```text
ionosctl logging-service pipeline key [flags]
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

Generate a new key for a logging pipeline, invalidating the old one. The key is used for authentication when sending logs.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -l, --location string      Location of the resource to operate on. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx (default "de/txl")
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string   The ID of the logging pipeline you want to generate a key for (required)
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service pipeline key --pipeline-id ID
```

