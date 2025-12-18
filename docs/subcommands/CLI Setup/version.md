---
description: "Show the current version"
---

# Version

## Usage

```text
ionosctl version [flags]
```

## Description

The `ionosctl version` command displays the current version of the ionosctl software and the latest Github release.

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --updates           Check for latest updates for CLI
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl version
```

