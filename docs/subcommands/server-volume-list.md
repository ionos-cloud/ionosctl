---
description: List attached Volumes from a Server
---

# ServerVolumeList

## Usage

```text
ionosctl server volume list [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `volume` command:

```text
[v vol]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Volumes attached to the Server.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose                see step by step process when running a command
```

## Examples

```text
ionosctl server volume list --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

