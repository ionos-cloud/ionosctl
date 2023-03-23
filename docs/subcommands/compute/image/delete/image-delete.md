---
description: Delete an image
---

# ImageDelete

## Usage

```text
ionosctl image delete [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Image.

Required values to run command:

* Image Id

## Options

```text
  -a, --all                Delete all non-public images
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -i, --image-id string    The unique Image Id (required)
      --no-headers         When using text output, don't print headers
  -t, --timeout int        Timeout option for Request for Image update [seconds] (default 60)
  -w, --wait-for-request   Wait for the Request for Image update to be executed
```

## Examples

```text
ionosctl image delete --image-id IMAGE_ID
```

