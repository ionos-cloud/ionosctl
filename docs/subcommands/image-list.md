---
description: List Images
---

# ImageList

## Usage

```text
ionosctl image list [flags]
```

## Aliases

For `image` command:
```text
[img]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of available public Images. Use flags to retrieve a list of sorted images by location, licence type, type or size.

## Options

```text
  -u, --api-url string        Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -h, --help                  help for list
      --licence-type string   The licence type of the Image
  -l, --location string       The location of the Image
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --size float32          The size of the Image
      --type string           The type of the Image
```

## Examples

```text
ionosctl image list

ionosctl image list --location us/las --type HDD --licence-type LINUX
```

