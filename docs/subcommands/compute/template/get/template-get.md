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
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --no-headers           When using text output, don't print headers
  -i, --template-id string   The unique Template Id (required)
```

## Examples

```text
ionosctl template get -i TEMPLATE_ID
```

