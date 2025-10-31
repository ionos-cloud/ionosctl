---
description: "Get a Snapshot"
---

# SnapshotGet

## Usage

```text
ionosctl snapshot get [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --snapshot-id string   The unique Snapshot Id (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl snapshot get --snapshot-id SNAPSHOT_ID
```

