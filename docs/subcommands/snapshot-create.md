---
description: Create a Snapshot of a Volume within the Virtual Data Center
---

# SnapshotCreate

## Usage

```text
ionosctl snapshot create [flags]
```

## Aliases

For `snapshot` command:
```text
[snap]
```

## Description

Use this command to create a Snapshot. Creation of Snapshots is performed from the perspective of the storage Volume. The name, description and licence type of the Snapshot can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Volume Id
* Name
* Licence Type

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -d, --description string     Description of the Snapshot
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
      --licence-type string    Licence Type of the Snapshot(required)
  -n, --name string            Name of the Snapshot(required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --sec-auth-protection    Enable secure authentication protection
  -t, --timeout int            Timeout option for Request for Snapshot creation [seconds] (default 60)
      --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Snapshot creation to be executed
```

## Examples

```text
ionosctl snapshot create --datacenter-id 451cc0c1-883a-44aa-9ae4-336c0c3eaa5d --volume-id 4acddd40-959f-4517-b628-dc24e37df942 --name testSnapshot
SnapshotId                             Name           LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   testSnapshot   LINUX         0
RequestId: fed5555a-ac00-41c8-abbe-cc53c8179716
Status: Command snapshot create has been successfully executed
```

