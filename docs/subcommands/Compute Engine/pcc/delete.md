---
description: "Delete a Cross-Connect"
---

# PccDelete

## Usage

```text
ionosctl pcc delete [flags]
```

## Aliases

For `pcc` command:

```text
[cc]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -a, --all                Delete all Cross-Connects.
  -u, --api-url string     Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -i, --pcc-id string      The unique Cross-Connect Id (required)
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for Cross-Connect deletion [seconds] (default 60)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for the Request for Cross-Connect deletion to be executed
```

## Examples

```text
ionosctl pcc delete --pcc-id PCC_ID --wait-for-request
```

