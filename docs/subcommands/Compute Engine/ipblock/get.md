---
description: "Get an IpBlock"
---

# IpblockGet

## Usage

```text
ionosctl ipblock get [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the attributes of a specific IpBlock.

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
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl ipblock get --ipblock-id IPBLOCK_ID
```

