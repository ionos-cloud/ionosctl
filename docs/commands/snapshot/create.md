---
description: Create a Snapshot of a Volume within the Virtual Data Center.
---

# Create

## Usage

```text
ionosctl snapshot create [flags]
```

## Description

Use this command to create a Snapshot in a specified Data Center. Creation of Snapshots is performed from the perspective of the storage volume. The name, description and licence type of the Snapshot can be set.

You can wait for the action to be executed using `--wait` option.

Required values to run command:
- Data Center Id
- Volume Id
- Snapshot Name

## Options

```text
  -u, --api-url string                 Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                   Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size])
  -c, --config string                  Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string           The unique Data Center Id [Required flag]
  -h, --help                           help for create
      --ignore-stdin                   Force command to execute without user input
  -o, --output string                  Desired output format [text|json] (default "text")
  -q, --quiet                          Quiet output
      --snapshot-description string    Description of the Snapshot
      --snapshot-licence-type string   Licence Type of the Snapshot
      --snapshot-name string           Name of the Snapshot [Required flag]
      --timeout int                    Timeout option for a Snapshot to be created [seconds] (default 60)
      --volume-id string               The unique Volume Id [Required flag]
      --wait                           Wait for Snapshot to be created
```

