---
description: Get a list of Peers from a Private Cross-Connect
---

# PccPeersList

## Usage

```text
ionosctl pcc peers list [flags]
```

## Description

Use this command to get a list of Peers from a Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -F, --format strings   Set of fields to be printed on output (default [LanId,LanName,DatacenterId,DatacenterName,Location])
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
      --pcc-id string    The unique Private Cross-Connect Id (required)
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl pcc peers list --pcc-id 4b9c6a43-a338-11eb-b70c-7ade62b52cc0 
LanId   LanName     DatacenterId                           DatacenterName   Location
1       testlan2    1ef56b51-98be-487e-925a-c9f3dfa4a076   test2            us/las
1       testlan1    95b7f7f0-a6f3-4fc9-8d06-018d2c1efc89   test1            us/las
```

