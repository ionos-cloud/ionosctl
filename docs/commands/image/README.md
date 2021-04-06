---
description: Image Operations
---

# Image

## Usage

```text
ionosctl image [command]
```

## Aliases

```text
[images img]
```

## Description

The sub-commands of `ionosctl image` allow you to see information about the Images available.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ImageId,Name,Location,Size,LicenceType,ImageType])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for image
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl image delete](delete.md) | Delete a private Image |
| [ionosctl image get](get.md) | Get a specified Image |
| [ionosctl image list](list.md) | List Images |

