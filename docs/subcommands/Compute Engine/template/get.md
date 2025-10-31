---
description: "Get a specified Template"
---

# TemplateGet

## Usage

```text
ionosctl template get [flags]
```

## Aliases

For `template` command:

```text
[tpl]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Template.

Required values to run command:

* Template Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TemplateId Name Cores RAM StorageSize] (default [TemplateId,Name,Cores,RAM,StorageSize])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --template-id string   The unique Template Id (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl template get -i TEMPLATE_ID
```

