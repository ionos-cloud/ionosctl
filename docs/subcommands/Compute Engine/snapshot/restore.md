---
description: "Restore a Snapshot onto a Volume"
---

# SnapshotRestore

## Usage

```text
ionosctl snapshot restore [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `restore` command:

```text
[r]
```

## Description

Use this command to restore a Snapshot onto a Volume. A Snapshot is created as just another image that can be used to create new Volumes or to restore an existing Volume.

Required values to run command:

* Datacenter Id
* Volume Id
* Snapshot Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud'|'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -i, --snapshot-id string     The unique Snapshot Id (required)
  -t, --timeout int            Timeout option for Request for Snapshot restore [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
      --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Snapshot restore to be executed
```

## Examples

```text
ionosctl snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --wait-for-request
```

