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

Use this command to delete a specified Virtual Data Center (VDC) from your account. This will remove all objects within the VDC and remove the VDC object itself. 

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location State Description Version Features SecAuthProtection] (default [DatacenterId,Name,Location,Features,State])
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
ionosctl datacenter delete --datacenter-id 8e543958-04f5-4872-bbf3-b28d46393ac7
Warning: Are you sure you want to delete data center (y/N) ? y
RequestId: 12547a71-9768-483b-8a8e-e03e58df6dc3
Status: Command datacenter delete has been successfully executed

ionosctl datacenter delete --datacenter-id ff279ffd-ac61-4e5d-ba5e-058296c77774 --force --wait-for-request
1.2s Waiting for request... DONE
RequestId: a2f71ef3-f81c-4b15-8f8f-5dfd1bdb3c26
Status: Command datacenter delete & wait have been successfully executed
```

