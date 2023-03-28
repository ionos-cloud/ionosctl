---
description: Get a specific attached CD-ROM from a Server
---

# ServerCdromGet

## Usage

```text
ionosctl server cdrom get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about an attached CD-ROM on Server.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -i, --cdrom-id string        The unique Cdrom Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl server cdrom get --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID
```

