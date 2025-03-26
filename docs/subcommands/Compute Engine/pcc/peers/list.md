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
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [LanId LanName DatacenterId DatacenterName Location] (default [LanId,LanName,DatacenterId,DatacenterName,Location])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
      --pcc-id string    The unique Cross-Connect Id (required)
  -q, --quiet            Quiet output
  -t, --timeout int      Timeout in seconds for polling the request (default 60)
  -v, --verbose          Print step-by-step process when running command
  -w, --wait             Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl pcc peers list --pcc-id PCC_ID
```

