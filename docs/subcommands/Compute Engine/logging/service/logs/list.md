---
description: "Retrieve logging pipeline logs"
---

# LoggingServiceLogsList

## Usage

```text
ionosctl logging-service logs list [flags]
```

## Aliases

For `logging-service` command:

```text
[log-svc]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve logging pipeline logs

## Options

```text
  -a, --all                  Use this flag to list all logging pipeline logs
  -u, --api-url string       Override default host url (default "logging.de-txl.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Tag Source Protocol Public Destinations] (default [Tag,Source,Protocol,Public,Destinations])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string   The ID of the logging pipeline you want to list logs for (required)
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service logs list --pipeline-id ID
```

