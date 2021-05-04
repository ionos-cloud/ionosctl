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
      --force             Force command to execute without user input
  -h, --help              help for get
      --image-id string   The unique Image Id (required)
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl image get --image-id 8fc5f591-338e-11eb-a681-1e659523cb7b 
ImageId                                Name                             Location   Size   LicenceType   ImageType
8fc5f591-338e-11eb-a681-1e659523cb7b   Ubuntu-19.10-server-2020-12-01   us/las     3      LINUX         HDD
```

