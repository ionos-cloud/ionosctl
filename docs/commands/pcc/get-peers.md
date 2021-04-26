---
description: Get a Private Cross-Connect Peers
---

# GetPeers

## Usage

```text
ionosctl pcc get-peers [flags]
```

## Description

Use this command to get a list of Peers from a Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for get-peers
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
      --pcc-id string    The unique Private Cross-Connect Id [Required flag]
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl pcc get-peers --pcc-id 4b9c6a43-a338-11eb-b70c-7ade62b52cc0 
LanId   LanName     DatacenterId                           DatacenterName   Location
1       testlan2    1ef56b51-98be-487e-925a-c9f3dfa4a076   test2            us/las
1       testlan1    95b7f7f0-a6f3-4fc9-8d06-018d2c1efc89   test1            us/las
```

