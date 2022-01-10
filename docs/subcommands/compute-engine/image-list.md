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

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name description location size cpuHotPlug cpuHotUnplug ramHotPlug ramHotUnplug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug discScsiHotPlug discScsiHotUnplug licenceType imageType public imageAliases cloudInit]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings       Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --image-alias string    Image Alias or part of Image Alias to sort Images by (deprecated)
      --latest int            Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed (deprecated)
      --licence-type string   The licence type of the Image (deprecated)
  -l, --location string       The location of the Image (deprecated)
  -M, --max-results int       The maximum number of elements to return
      --order-by string       Limits results to those containing a matching value for a specific property
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --type string           The type of the Image (deprecated)
  -v, --verbose               Print step-by-step process when running command
```

## Examples

```text
ionosctl image list

ionosctl image list --location us/las --type HDD --licence-type LINUX
```

