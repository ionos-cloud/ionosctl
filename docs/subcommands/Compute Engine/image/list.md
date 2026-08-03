---
description: "List public and private Images"
---

# ImageList

## Usage

```text
ionosctl compute image list [flags]
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

List the Images your contract can see: IONOS-provided PUBLIC images (OS distributions and ISOs) and your own PRIVATE images (uploaded or snapshotted). The Public column distinguishes the two; the ImageAliases column shows the stable names (e.g. "ubuntu:latest") you can pass instead of a UUID when creating a volume.

Because public images are replicated per location, the same OS appears once per region (see the Location column, e.g. de/fra, gb/lhr). Narrow the result with server-side filters.

You can filter the results using the `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`. To list only your own images, filter on `public=false`.
Available Filters:
* filter by property: [name description location size cpuHotPlug cpuHotUnplug ramHotPlug ramHotUnplug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug discScsiHotPlug discScsiHotUnplug exposeSerial requireLegacyBios licenceType applicationType imageType public imageAliases requiredFeatures cloudInit]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string        Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [ImageId Name ImageAliases Location LicenceType ImageType CloudInit CreatedDate Size Description Public CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType RequiredFeatures]
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int             Level of detail for response objects (default 1)
  -F, --filters strings       Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --image-alias string    Client-side filter keeping only images whose image-alias contains this substring. An image-alias (e.g. "ubuntu:20.04", "debian:latest") is a stable human-friendly name you can use in place of a UUID when creating a volume (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
      --latest int            Client-side: keep only the N most recently created images (by createdDate, newest first). 0 (default) keeps all. Prefer --order-by createdDate --max-results N for server-side ordering (DEPRECATED: Use --filters --order-by --max-results options instead!)
      --licence-type string   Client-side filter by OS licence type (LINUX, RHEL, WINDOWS, WINDOWS2016/2019/2022/2025, UNKNOWN, OTHER). This is how the platform bills and configures the guest OS (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
      --limit int             Maximum number of items to return per request (default 50)
  -l, --location string       Client-side filter by the physical location an image lives in, e.g. de/fra, de/txl, gb/lhr. Public images are replicated per location, so the same OS appears once per region (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
      --no-headers            Don't print table headers when table output is used
      --offset int            Number of items to skip before starting to collect the results
      --order-by string       Property to order the results by
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
  -t, --timeout int           Timeout in seconds for --wait and other wait operations (default 600)
      --type string           Client-side filter by image type: HDD (a bootable disk image) or CDROM (an ISO you can attach as a virtual optical drive) (DEPRECATED: incompatible with --max-results. Use --filters --order-by --max-results options instead!)
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                  Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# List every image visible to your contract (public + private)
ionosctl compute image list

# List only your own uploaded/snapshotted images
ionosctl compute image list --filters public=false

# Find public Ubuntu HDD images in Frankfurt via server-side filters
ionosctl compute image list --filters public=true,imageAliases=ubuntu:latest,location=de/fra
```

