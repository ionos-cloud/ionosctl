---
description: Create a Snapshot
---

# Create

## Usage

```text
ionosctl snapshot create [flags]
```

## Description

Use this command to create a Server in a specified Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the action to be executed using `--wait` option.

Required values to run command:
- Data Center Id

## Options

```text
  -u, --api-url string                Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                  Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size])
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string          The unique Snapshot Id [Required flag]
  -h, --help                          help for create
      --ignore-stdin                  Force command to execute without user input
  -o, --output string                 Desired output format [text|json] (default "text")
  -q, --quiet                         Quiet output
      --snapshot-description string   CPU Family for the Server
      --snapshot-id string            The unique Snapshot Id [Required flag]
      --snapshot-licencetype string   CPU Family for the Server
      --snapshot-name string          Name of the Server
      --timeout int                   Timeout option [seconds] (default 60)
      --volume-id string              The unique Snapshot Id [Required flag]
      --wait                          Wait for Server to be created
```

