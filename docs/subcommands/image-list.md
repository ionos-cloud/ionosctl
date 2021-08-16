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

Use this command to get a full list of available public Images.

Use flags to retrieve a list of Images:

* sorting by location, using `ionosctl image list --location LOCATION_ID`
* sorting by licence type, using `ionosctl image list --licence-type LICENCE_TYPE`
* sorting by Image type, using `ionosctl image list --type IMAGE_TYPE`
* sorting by Image alias, using `ionosctl image list --image-alias IMAGE_ALIAS`; image alias can be either the Image alias `--image-alias ubuntu:latest` or part of Image alias e.g. `--image-alias latest`
* sorting by the time the Image was created, starting from now in descending order, take the first N Images, using `ionosctl image list --latest N`
* sorting by multiple of above options, using `ionosctl image list --type IMAGE_TYPE --location LOCATION_ID --latest N`

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -h, --help                  help for list
      --image-alias string    Image Alias or part of Image Alias to sort Images by
      --latest int            Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed
      --licence-type string   The licence type of the Image
  -l, --location string       The location of the Image
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --type string           The type of the Image
  -v, --verbose               see step by step process when running a command
```

## Examples

```text
ionosctl image list

ionosctl image list --location us/las --type HDD --licence-type LINUX
```

