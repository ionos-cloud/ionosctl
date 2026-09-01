---
description: "Restore a Snapshot onto a Volume (destructive overwrite)"
---

# SnapshotRestore

## Usage

```text
ionosctl compute snapshot restore [flags]
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

Use this command to write the contents of a Snapshot back onto an existing target Volume, rolling that Volume to the point in time the Snapshot captured.

This is a DESTRUCTIVE, in-place overwrite: the target Volume's current data is replaced by the Snapshot image, so anything written since is lost. You are asked to confirm (use --force to skip the prompt). The Snapshot and the target Volume must be at the same LOCATION, and the target Volume should be at least as large as the Snapshot.

The target Volume is identified by its Data Center Id + Volume Id; it does NOT have to be the Volume the Snapshot was originally taken from - any compatible Volume at the same location works. To spin up a brand-new Volume from a Snapshot instead of overwriting one, pass the Snapshot Id as the image when creating a Volume rather than using restore.

Use `--wait` (`-w`) to wait until the restore completes and the Volume is AVAILABLE.

Required values to run command:

* Datacenter Id
* Volume Id
* Snapshot Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [SnapshotId Name LicenceType Size State]
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
  -i, --snapshot-id string     The unique Snapshot Id (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The unique Volume Id (required)
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Roll a volume back to a snapshot (prompts for confirmation)
ionosctl compute snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --wait

# Advanced: restore onto a different target volume, no prompt, in a script
ionosctl compute snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id OTHER_VOLUME_ID --force --wait
```

