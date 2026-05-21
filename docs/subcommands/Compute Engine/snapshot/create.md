---
description: "Create a Snapshot of a Volume within the Virtual Data Center"
---

# SnapshotCreate

## Usage

```text
ionosctl compute snapshot create [flags]
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

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [SnapshotId Name LicenceType Size State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --description string     Description of the Snapshot
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --licence-type string    Licence Type of the Snapshot. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "LINUX")
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            Name of the Snapshot (default "Unnamed Snapshot")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --sec-auth-protection    Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The unique Volume Id (required)
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name NAME
```

