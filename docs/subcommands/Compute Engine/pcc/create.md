---
description: "Create a Cross-Connect"
---

# PccCreate

## Usage

```text
ionosctl compute pcc create [flags]
```

## Aliases

For `pcc` command:

```text
[cc]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Cross-Connect (Private Cross-Connect).

This creates the empty Cross-Connect object only; no LAN is attached yet. After creation, attach LANs by running `ionosctl compute lan update --datacenter-id <DC_ID> --lan-id <LAN_ID> --pcc-id <PCC_ID>` for each private LAN you want to peer. All LANs you attach must be in the same contract and region and must have non-overlapping private IP ranges.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [PccId Name Description State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -d, --description string   A human-readable description for the Cross-Connect
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -n, --name string          The name for the Cross-Connect. Purely a display label; does not affect connectivity (default "Unnamed PrivateCrossConnect")
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a Cross-Connect with a name and description
ionosctl compute pcc create --name "prod-interconnect" --description "Bridges DB and app VDCs"

# Create and wait for the request to complete
ionosctl compute pcc create --name "prod-interconnect" --wait
```

