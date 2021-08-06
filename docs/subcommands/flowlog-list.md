---
description: List FlowLogs
---

# FlowlogList

## Usage

```text
ionosctl flowlog list [flags]
```

## Aliases

For `flowlog` command:

```text
[fl]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of FlowLogs from a specified NIC from a Server.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl flowlog list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

