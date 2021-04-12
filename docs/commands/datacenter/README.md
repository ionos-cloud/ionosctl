---
description: Data Center Operations
---

# DataCenter

## Usage

```text
ionosctl datacenter [command]
```

## Aliases

```text
[dc]
```

## Description

The sub-commands of `ionosctl datacenter` allow you to create, list, get, update and delete Data Centers.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for datacenter
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl datacenter add-label](add-label.md) | Add a Label to a Data Center |
| [ionosctl datacenter create](create.md) | Create a Data Center |
| [ionosctl datacenter delete](delete.md) | Delete a Data Center |
| [ionosctl datacenter get](get.md) | Get a Data Center |
| [ionosctl datacenter get-label](get-label.md) | Get a Label from a Data Center |
| [ionosctl datacenter list](list.md) | List Data Centers |
| [ionosctl datacenter list-labels](list-labels.md) | List Labels from a Data Center |
| [ionosctl datacenter remove-label](remove-label.md) | Remove a Label from a Data Center |
| [ionosctl datacenter update](update.md) | Update a Data Center |

