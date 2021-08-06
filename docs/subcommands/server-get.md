---
description: Get a Server
---

# ServerGet

## Usage

```text
ionosctl server get [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Server from a Virtual Data Center. You can also wait for Server to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ServerId Name AvailabilityZone Cores Ram CpuFamily VmState State TemplateId Type] (default [ServerId,Name,Type,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for waiting for Server to be in AVAILABLE state [seconds] (default 60)
  -W, --wait-for-state         Wait for specified Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server get --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

