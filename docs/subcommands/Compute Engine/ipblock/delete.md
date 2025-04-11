---
description: "Delete an IpBlock"
---

# IpblockDelete

## Usage

```text
ionosctl ipblock delete [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified IpBlock.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -a, --all                 Delete all the IpBlocks.
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -i, --ipblock-id string   The unique IpBlock Id (required)
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout duration    Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl ipblock delete --ipblock-id IPBLOCK_ID --wait-for-request
```

