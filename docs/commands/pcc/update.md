---
description: Update a Private Cross-Connect
---

# Update

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
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [PccId,Name,Description])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                    Force command to execute without user input
  -h, --help                     help for update
  -o, --output string            Desired output format [text|json] (default "text")
      --pcc-description string   The description for the Private Cross-Connect
      --pcc-id string            The unique Private Cross-Connect Id (required)
      --pcc-name string          The name for the Private Cross-Connect
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for Private Cross-Connect to be updated [seconds] (default 60)
      --wait                     Wait for Private Cross-Connect to be updated
```

## Examples

```text
ionosctl pcc update --pcc-id 4b9c6a43-a338-11eb-b70c-7ade62b52cc0 --pcc-description test
PccId                                  Name   Description
4b9c6a43-a338-11eb-b70c-7ade62b52cc0   test   test
RequestId: 81525f2d-cc91-4c55-84b8-07fac9a47e35
Status: Command pcc update has been successfully executed
```

