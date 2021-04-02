---
description: Get a specified Image
---

# Get

## Usage

```text
ionosctl image get [flags]
```

## Description

Use this command to get information about a specified Image.

Required values to run command:

* Image Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [ImageId,Name,Location,Size,LicenceType,ImageType])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help              help for get
      --ignore-stdin      Force command to execute without user input
      --image-id string   The unique Image Id [Required flag]
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

