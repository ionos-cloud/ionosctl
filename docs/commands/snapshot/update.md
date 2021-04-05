---
description: Update a Snapshot.
---

# Update

## Usage

```text
ionosctl snapshot update [flags]
```

## Description

Use this command to update a specified Snapshot.

You can wait for the action to be executed using `--wait` option.

Required values to run command:
- Snapshot Id

## Options

```text
  -u, --api-url string                 Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                   Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string                  Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                           help for update
      --ignore-stdin                   Force command to execute without user input
  -o, --output string                  Desired output format [text|json] (default "text")
  -q, --quiet                          Quiet output
      --snapshot-description string    Description of the Snapshot
      --snapshot-id string             The unique Snapshot Id [Required flag]
      --snapshot-licence-type string   Licence Type of the Snapshot
      --snapshot-name string           Name of the Snapshot
      --snapshot-size float32          Size of the Snapshot (default 10)
      --timeout int                    Timeout option for a Snapshot to be created [seconds] (default 60)
      --wait                           Wait for Snapshot to be created
```

## Examples

```text
ionosctl snapshot update --snapshot-id dc688daf-8e54-4db8-ac4a-487ad5a34e9c --snapshot-name test
SnapshotId                             Name   LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   test   LINUX         10
RequestId: 3540e9be-ed35-41c0-83d9-923882bfa9bd
Status: Command snapshot update has been successfully executed
```

