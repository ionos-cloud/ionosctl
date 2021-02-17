---
description: Force a hard reboot of a Server
---

# Reboot

## Usage

```text
ionosctl server reboot [flags]
```

## Description

Use this command to force a hard reboot of the Server. Do not use this method if you want to gracefully reboot the machine. This is the equivalent of powering off the machine and turning it back on.

You can wait for the action to be executed using `--wait` option.
You can force the command to execute without user input using `--ignore-stdin` option.

Required values to run command:
- Data Center Id
- Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for reboot
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
  -v, --verbose                Enable verbose output
```

## Examples

```text
ionosctl server reset --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa
⚠ Warning: Are you sure you want to reboot server (y/N) ? y
✔ RequestId: e6720605-2fa4-46d9-be74-42b733eb1128
✔ Status: Command server reset has been successfully executed
```

## See also

* [ionosctl server](./)

