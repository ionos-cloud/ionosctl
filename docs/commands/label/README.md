---
description: Label Operations
---

# Label

## Usage

```text
ionosctl label [command]
```

## Description

The sub-commands of `ionosctl label` allow you to get, list, create, delete Labels from a resource. For each resource that supports labelling: Data Center, Server, Volume, IPBlock, Snapshot - commands to manage Labels are available. Example: `ionosctl <resource> add-label`.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for label
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl label get](get.md) | Get a Label |
| [ionosctl label list](list.md) | List Labels from all Resources |

