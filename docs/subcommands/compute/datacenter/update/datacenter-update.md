---
description: Update a Data Center
---

# DatacenterUpdate

## Usage

```text
ionosctl datacenter update [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to change a Virtual Data Center's name, description.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string     Description of the Data Center
  -n, --name string            Name of the Data Center
  -t, --timeout int            Timeout option for Request for Data Center update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Data Center update to be executed
```

## Examples

```text
ionosctl datacenter update --datacenter-id DATACENTER_ID --description DESCRIPTION --cols "DatacenterId,Description"
```

