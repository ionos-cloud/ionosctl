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
  -D, --depth int32           Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings       Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
      --image-alias string    Image Alias or part of Image Alias to sort Images by (deprecated)
      --latest int            Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed (deprecated)
      --licence-type string   The licence type of the Image (deprecated)
  -l, --location string       The location of the Image (deprecated)
  -M, --max-results int32     The maximum number of elements to return
      --no-headers            When using text output, don't print headers
      --order-by string       Limits results to those containing a matching value for a specific property
      --type string           The type of the Image (deprecated)
```

## Examples

```text
ionosctl image list

ionosctl image list --location us/las --type HDD --licence-type LINUX
```

