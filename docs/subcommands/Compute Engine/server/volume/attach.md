---
description: "Attach a Volume to a Server"
---

# ServerVolumeAttach

## Usage

```text
ionosctl server volume attach [flags]
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

For `attach` command:

```text
[a]
```

## Description

Use this command to attach a pre-existing Volume to a Server.

You can wait for the Request to be executed using `--wait-for-request` option.

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
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Volume attachment [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -i, --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Volume attachment to be executed
```

## Examples

```text
ionosctl server volume attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID
```

