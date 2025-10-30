---
description: "Update a Cross-Connect"
---

# PccUpdate

## Usage

```text
ionosctl pcc update [flags]
```

## Aliases

For `pcc` command:

```text
[cc]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update details about a specific Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string   The description for the Cross-Connect
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -n, --name string          The name for the Cross-Connect
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pcc-id string        The unique Cross-Connect Id (required)
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Cross-Connect update [seconds] (default 60)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request     Wait for the Request for Cross-Connect update to be executed
```

## Examples

```text
ionosctl pcc update --pcc-id PCC_ID --description DESCRIPTION
```

