---
description: "Delete a Snapshot"
---

# SnapshotDelete

## Usage

```text
ionosctl snapshot delete [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete the specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -a, --all                  Delete all the Snapshots.
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print column headers
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --snapshot-id string   The unique Snapshot Id (required)
  -t, --timeout int          Timeout option for Request for Snapshot deletion [seconds] (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait-for-request     Wait for the Request for Snapshot deletion to be executed
```

## Examples

```text
ionosctl snapshot delete --snapshot-id SNAPSHOT_ID --wait-for-request
```

