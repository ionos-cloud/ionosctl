---
description: Update a Data Center
---

# Update

## Usage

```text
ionosctl datacenter update [flags]
```

## Description

Use this command to change a Data Center's name, description. 

You can wait for the action to be executed using `--wait` option.

Required values to run command:
- Data Center Id

## Options

```text
  -u, --api-url string                  Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                    Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-description string   Description of the Data Center
      --datacenter-id string            The unique Data Center Id [Required flag]
      --datacenter-name string          Name of the Data Center
  -h, --help                            help for update
      --ignore-stdin                    Force command to execute without user input
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
      --timeout int                     Timeout option [seconds] (default 60)
  -v, --verbose                         Enable verbose output
      --wait                            Wait for Data Center to be updated
```

## Examples

```text
ionosctl datacenter update --datacenter-id 8e543958-04f5-4872-bbf3-b28d46393ac7 --datacenter-description demoDescription --cols "DatacenterId,Description"
DatacenterId                           Description
8e543958-04f5-4872-bbf3-b28d46393ac7   demoDescription
✔ RequestId: 46af6915-9003-4f11-a1fe-bab1eac9bccc
✔ Status: Command datacenter update has been successfully executed
```

## See also

* [ionosctl datacenter](./)

