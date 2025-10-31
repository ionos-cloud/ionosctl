---
description: "Get an attached Volume from a Server"
---

# ServerVolumeGet

## Usage

```text
ionosctl server volume get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about an attached Volume on Server.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              pagination limit: Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -i, --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl server volume get --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID
```

