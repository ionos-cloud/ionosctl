---
description: List images
---

# List

## Usage

```text
ionosctl image list [flags]
```

## Description

Use this command to get a list of available public images. Use flags to retrieve a list of sorted images by location, licence type, type or size.

## Options

```text
  -u, --api-url string              Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                Columns to be printed in the standard output (default [ImageId,Name,Location,Size,LicenceType,ImageType])
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                        help for list
      --ignore-stdin                Force command to execute without user input
      --image-licence-type string   
      --image-location string       
      --image-size float32          
      --image-type string           
  -o, --output string               Desired output format [text|json] (default "text")
  -q, --quiet                       Quiet output
```

