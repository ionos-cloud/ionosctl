---
description: "List attached CD-ROMs from a Server"
---

# ServerCdromList

## Usage

```text
ionosctl server cdrom list [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `cdrom` command:

```text
[cd]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of CD-ROMs attached to the Server.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name description location size cpuHotPlug cpuHotUnplug ramHotPlug ramHotUnplug nicHotPlug nicHotUnplug discVirtioHotPlug discVirtioHotUnplug discScsiHotPlug discScsiHotUnplug exposeSerial requireLegacyBios licenceType applicationType imageType public imageAliases cloudInit]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl server cdrom list --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

