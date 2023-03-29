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
[d dc vdc]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Virtual Data Center from your account. This will remove all objects within the VDC and remove the VDC object itself.

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    Delete all the Datacenters.
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -t, --timeout int            Timeout option for Request for Data Center deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Data Center deletion
```

## Examples

```text
ionosctl datacenter delete --datacenter-id DATACENTER_ID
ionosctl datacenter delete --datacenter-id DATACENTER_ID --force --wait-for-request
```

