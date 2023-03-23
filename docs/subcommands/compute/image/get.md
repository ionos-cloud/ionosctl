---
description: Get a specified Image
---

# ImageGet

## Usage

```text
ionosctl image get [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Image.

Required values to run command:

* Image Id

## Options

```text
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -i, --image-id string   The unique Image Id (required)
      --no-headers        When using text output, don't print headers
```

## Examples

```text
ionosctl image get --image-id IMAGE_ID
```

