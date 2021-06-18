---
description: Get a list of Peers from a Private Cross-Connect
---

# PccPeersList

## Usage

```text
ionosctl pcc peers list [flags]
```

## Aliases

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Peers from a Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [LanId LanName DatacenterId DatacenterName Location] (default [LanId,LanName,DatacenterId,DatacenterName,Location])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
      --pcc-id string    The unique Private Cross-Connect Id (required)
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl pcc peers list --pcc-id PCC_ID
```

