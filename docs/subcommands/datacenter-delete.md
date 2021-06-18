---
description: Delete a Data Center
---

# DatacenterDelete

## Usage

```text
ionosctl datacenter delete [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Virtual Data Center \(VDC\) from your account. This will remove all objects within the VDC and remove the VDC object itself.

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection] (default [DatacenterId,Name,Location,CpuFamily,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for Data Center deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Data Center deletion
```

## Examples

```text
ionosctl datacenter delete --datacenter-id DATACENTER_ID

ionosctl datacenter delete --datacenter-id DATACENTER_ID --force --wait-for-request
```

