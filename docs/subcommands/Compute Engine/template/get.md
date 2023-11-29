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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TemplateId Name Cores Ram StorageSize] (default [TemplateId,Name,Cores,Ram,StorageSize])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --template-id string   The unique Template Id (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl template get -i TEMPLATE_ID
```

