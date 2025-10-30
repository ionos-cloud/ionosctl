---
description: "Update an IpBlock"
---

# IpblockUpdate

## Usage

```text
ionosctl ipblock update [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the properties of an existing IpBlock.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -i, --ipblock-id string   The unique IpBlock Id (required)
  -n, --name string         Name of the IpBlock
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for Request for IpBlock update [seconds] (default 60)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request    Wait for the Request for IpBlock update to be executed
```

## Examples

```text
ionosctl ipblock update --ipblock-id IPBLOCK_ID --ipblock-name NAME
```

