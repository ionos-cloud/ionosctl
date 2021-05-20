---
description: Get a specified Image
---

# ImageGet

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
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -F, --format strings    Set of fields to be printed on output (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
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

