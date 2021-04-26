---
description: Delete a private Image
---

# Delete

## Usage

```text
ionosctl image delete [flags]
```

## Description

Use this command to delete the specified private image. This only applies to private images that you have uploaded.

Required values to run command:

* Image Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [ImageId,Name,Location,Size,LicenceType,ImageType])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help              help for delete
      --force      Force command to execute without user input
      --image-id string   The unique Image Id [Required flag]
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

