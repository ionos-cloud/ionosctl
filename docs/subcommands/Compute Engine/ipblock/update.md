---
description: "Update an IpBlock"
---

# IpblockUpdate

## Usage

```text
ionosctl compute ipblock update [flags]
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

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -i, --ipblock-id string   The unique IpBlock Id (required)
      --limit int           Maximum number of items to return per request (default 50)
  -n, --name string         Name of the IpBlock
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute ipblock update --ipblock-id IPBLOCK_ID --ipblock-name NAME
```

