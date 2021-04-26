---
description: Restore a Snapshot onto a Volume
---

# Restore

## Usage

```text
ionosctl snapshot restore [flags]
```

## Description

Use this command to restore a Snapshot onto a Volume. A Snapshot is created as just another image that can be used to create new Volumes or to restore an existing Volume.

Required values to run command:

* Datacenter Id
* Volume Id
* Snapshot Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
      --force                  Force command to execute without user input
  -h, --help                   help for restore
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --snapshot-id string     The unique Snapshot Id [Required flag]
      --timeout int            Timeout option for a Snapshot to be restored [seconds] (default 60)
      --volume-id string       The unique Volume Id [Required flag]
      --wait                   Wait for Snapshot to be restored
```

## Examples

```text
ionosctl snapshot restore --snapshot-id dc688daf-8e54-4db8-ac4a-487ad5a34e9c --datacenter-id 451cc0c1-883a-44aa-9ae4-336c0c3eaa5d --volume-id 4acddd40-959f-4517-b628-dc24e37df942 --wait 
Warning: Are you sure you want to restore snapshot (y/N) ? 
y
RequestId: 21ca5546-9314-4cd5-8832-6029638b1237
Status: Command snapshot restore and request have been successfully executed
```

