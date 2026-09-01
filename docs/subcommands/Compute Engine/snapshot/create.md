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

Use this command to take a Snapshot of a storage Volume. A Snapshot is created from the perspective of a specific Volume, so both the Volume Id and the Data Center Id that Volume lives in are required.

The Snapshot captures the FULL provisioned capacity of the Volume (not just the used space) as a point-in-time image, and is stored independently at the Volume's LOCATION. From then on it can be restored onto a Volume (`snapshot restore`) or used as a boot image when creating new Volumes - but only within that same location and your contract.

For a consistent image, prefer taking the Snapshot while the source Volume's Server is powered off (or with in-guest I/O quiesced); a Snapshot of a live, busy filesystem may be crash-consistent only.

Use `--wait` (`-w`) to wait for the Snapshot to reach AVAILABLE state before the command returns.

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
  -d, --description string     Free-form notes about the Snapshot, e.g. why or when it was taken
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --licence-type string    The operating-system licence recorded on the Snapshot. This carries over to Volumes created/restored from it and affects OS licensing (notably WINDOWS variants are billed). Set it to match the source Volume's OS. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "LINUX")
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            A human-friendly label for the Snapshot; shown in listings. Does not have to be unique (default "Unnamed Snapshot")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --sec-auth-protection    Protect the Snapshot with secure authentication: when true, deleting or restoring it requires the Contract Owner or a re-authenticated user, guarding against accidental loss. E.g.: --sec-auth-protection=true, --sec-auth-protection=false
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The unique Volume Id (required)
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Take a basic snapshot of a volume
ionosctl compute snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name "pre-upgrade baseline"

# Advanced: name it, note the OS licence, and require Contract Owner / re-authentication to delete or restore it
ionosctl compute snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name "prod-db golden" --description "before schema migration" --licence-type LINUX --sec-auth-protection=true --wait
```

