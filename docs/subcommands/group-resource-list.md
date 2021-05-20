---
description: List Resources from a Group
---

# GroupResourceList

## Usage

```text
ionosctl group resource list [flags]
```

## Description

Use this command to get a list of Resources assigned to a Group. To see more details about existing Resources, use `ionosctl resource` commands.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [ResourceId,Name,SecAuthProtection,Type,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              help for list
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl group resource list --group-id 45ba215b-6897-40b6-879c-cbadb527cefd 
ResourceId                             Name   SecAuthProtection   Type
aa8e07a2-287a-4b45-b5e9-94761750a53c   test   false               datacenter
```

