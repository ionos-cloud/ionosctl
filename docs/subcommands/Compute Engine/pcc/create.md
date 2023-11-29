---
description: "Create a Cross-Connect"
---

# PccCreate

## Usage

```text
ionosctl pcc create [flags]
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

Use this command to create a Cross-Connect. You can specify the name and the description for the Cross-Connect.

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string   The description for the Cross-Connect
  -h, --help                 Print usage
  -n, --name string          The name for the Cross-Connect (default "Unnamed PrivateCrossConnect")
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Cross-Connect creation [seconds] (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait-for-request     Wait for the Request for Cross-Connect creation to be executed
```

## Examples

```text
ionosctl pcc create --name NAME --description DESCRIPTION --wait-for-request
```

