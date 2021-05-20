---
description: Update a Private Cross-Connect
---

# PccUpdate

## Usage

```text
ionosctl pcc update [flags]
```

## Description

Use this command to update details about a specific Private Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --description string   The description for the Private Cross-Connect
  -f, --force                Force command to execute without user input
  -F, --format strings       Set of fields to be printed on output (default [PccId,Name,Description,State])
  -h, --help                 help for update
  -n, --name string          The name for the Private Cross-Connect
  -o, --output string        Desired output format [text|json] (default "text")
      --pcc-id string        The unique Private Cross-Connect Id (required)
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Private Cross-Connect update [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Private Cross-Connect update to be executed
```

## Examples

```text
ionosctl pcc update --pcc-id 4b9c6a43-a338-11eb-b70c-7ade62b52cc0 --pcc-description test
PccId                                  Name   Description
4b9c6a43-a338-11eb-b70c-7ade62b52cc0   test   test
RequestId: 81525f2d-cc91-4c55-84b8-07fac9a47e35
Status: Command pcc update has been successfully executed
```

