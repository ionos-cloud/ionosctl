---
description: Get a Token from a Server
---

# ServerTokenGet

## Usage

```text
ionosctl server token get [flags]
```

## Aliases

For `server` command:
```text
[s svr]
```

For `token` command:
```text
[t]
```

For `get` command:
```text
[g]
```

## Description

Use this command to get the Server's jwToken.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ServerId Name AvailabilityZone Cores Ram CpuFamily VmState State TemplateId Type] (default [ServerId,Name,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl server token get --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

