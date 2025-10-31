---
description: "Get a list of Peers from a Cross-Connect"
---

# PccPeersList

## Usage

```text
ionosctl pcc peers list [flags]
```

## Aliases

For `pcc` command:

```text
[cc]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Peers from a Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [LanId LanName DatacenterId DatacenterName Location] (default [LanId,LanName,DatacenterId,DatacenterName,Location])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
      --pcc-id string    The unique Cross-Connect Id (required)
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl pcc peers list --pcc-id PCC_ID
```

