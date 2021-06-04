---
description: Resume a Cube Server
---

# ServerResume

## Usage

```text
ionosctl server resume [flags]
```

## Aliases

For `server` command:
```text
[s svr]
```

For `resume` command:
```text
[res]
```

## Description

Use this command to resume a Cube Server. The operation can only be applied to suspended Cube Servers.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

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
  -h, --help                   help for resume
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server resume [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server resume to be executed
```

## Examples

```text
ionosctl server resume --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

