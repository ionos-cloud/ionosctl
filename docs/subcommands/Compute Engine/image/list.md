---
description: "List Images"
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
* filter by property: [name description location size cpuHotPlug cpuHotUnplug ramHotPlug ramHotUnplug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug discScsiHotPlug discScsiHotUnplug exposeSerial requireLegacyBios licenceType applicationType imageType public imageAliases cloudInit]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string        Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int             Level of detail for response objects (default 1)
  -F, --filters strings       Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --image-alias string    Image Alias or part of Image Alias to sort Images by (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
      --latest int            Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed (DEPRECATED: Use --filters --order-by --max-results options instead!)
      --licence-type string   The licence type of the Image (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
      --limit int             Maximum number of items to return per request (default 50)
  -l, --location string       The location of the Image (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
      --no-headers            Don't print table headers when table output is used
      --offset int            Number of items to skip before starting to collect the results
      --order-by string       Property to order the results by
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
      --type string           The type of the Image (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl image list

ionosctl image list --location us/las --type HDD --licence-type LINUX
```

