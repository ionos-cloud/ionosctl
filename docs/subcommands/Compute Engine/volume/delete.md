---
description: "Delete a Volume"
---

# VolumeDelete

## Usage

```text
ionosctl compute volume delete [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete specified Volume. This will result in the Volume being removed from your Virtual Data Center. Please use this with caution!

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -a, --all                    Delete all Volumes from a virtual Datacenter.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for Volume deletion [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -i, --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Volume deletion to be executed
```

## Examples

```text
ionosctl compute volume delete --datacenter-id DATACENTER_ID --volume-id VOLUME_ID
```

