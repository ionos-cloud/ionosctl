---
description: Volume Operations
---

# Volume

## Usage

```text
ionosctl volume [command]
```

## Description

The sub-commands of `ionosctl volume` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes or attaching, detaching a Volume from a Server.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for volume
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Enable verbose output
```

## See also

* [ionosctl](../)
* [ionosctl volume attach](attach/)
* [ionosctl volume create](create.md)
* [ionosctl volume delete](delete.md)
* [ionosctl volume detach](detach.md)
* [ionosctl volume get](get.md)
* [ionosctl volume list](list.md)
* [ionosctl volume update](update.md)

