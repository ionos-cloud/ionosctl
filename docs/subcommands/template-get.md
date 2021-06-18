---
description: Get a specified Template
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
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TemplateId Name Cores Ram StorageSize] (default [TemplateId,Name,Cores,Ram,StorageSize])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for get
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --template-id string   The unique Template Id (required)
```

## Examples

```text
ionosctl template get -i TEMPLATE_ID
```

