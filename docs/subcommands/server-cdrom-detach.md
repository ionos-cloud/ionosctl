---
description: Detach a CD-ROM from a Server
---

# ServerCdromDetach

## Usage

```text
ionosctl server cdrom detach [flags]
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

For `detach` command:
```text
[d]
```

## Description

This will detach the CD-ROM from the Server.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -i, --cdrom-id string        The unique Cdrom Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for detach
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for CD-ROM detachment [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for CD-ROM detachment to be executed
```

## Examples

```text
ionosctl server cdrom detach --datacenter-id 4fd7996d-2b08-4c04-9c47-d9d884ee179a --server-id f7438b0c-2f36-4bec-892f-af027930b81e --cdrom-id 80c63662-49a0-11ea-94e0-525400f64d8d --wait-for-request --force 
15s Waiting for request. DONE                                                                                                                                                                              
RequestId: f6eb5b4a-7eb9-4515-872a-dd5acc296968
Status: Command cdrom detach & wait have been successfully executed
```

