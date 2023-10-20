---
description: "Suspend a Cube Server"
---

# ServerSuspend

## Usage

```text
ionosctl server suspend [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

## Description

Use this command to suspend a Cube Server. The operation can only be applied to Cube Servers. Note: The virtual machine will not be deleted.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ServerId DatacenterId Name AvailabilityZone Cores Ram CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId] (default [ServerId,Name,Type,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server suspend [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for Server suspend to be executed
```

## Examples

```text
ionosctl server suspend --datacenter-id DATACENTER_ID -i SERVER_ID
```

