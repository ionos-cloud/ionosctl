---
description: Delete a Private Cross-Connect
---

# PccDelete

## Usage

```text
ionosctl pcc delete [flags]
```

## Aliases

For `delete` command:
```text
[d]
```

## Description

Use this command to delete a Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               help for delete
  -o, --output string      Desired output format [text|json] (default "text")
  -i, --pcc-id string      The unique Private Cross-Connect Id (required)
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for Private Cross-Connect deletion [seconds] (default 60)
  -w, --wait-for-request   Wait for the Request for Private Cross-Connect deletion to be executed
```

## Examples

```text
ionosctl pcc delete --pcc-id e2337b40-52d9-48d2-bcbc-41c5abc29d11 --wait-for-request
Warning: Are you sure you want to delete private cross-connect (y/N) ? 
y
RequestId: 7fa56e7f-1d63-4c5f-a7ea-eec6a015282a
Status: Command pcc delete & wait have been successfully executed
```

