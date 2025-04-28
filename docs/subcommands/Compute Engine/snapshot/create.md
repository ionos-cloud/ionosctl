---
description: "Create a Snapshot of a Volume within the Virtual Data Center"
---

# SnapshotCreate

## Usage

```text
ionosctl snapshot create [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Snapshot. Creation of Snapshots is performed from the perspective of the storage Volume. The name, description and licence type of the Snapshot can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string     Description of the Snapshot
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --licence-type string    Licence Type of the Snapshot. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "LINUX")
  -n, --name string            Name of the Snapshot (default "Unnamed Snapshot")
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --sec-auth-protection    Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false
  -t, --timeout int            Timeout option for Request for Snapshot creation [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
      --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Snapshot creation to be executed
```

## Examples

```text
ionosctl snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name NAME
```

