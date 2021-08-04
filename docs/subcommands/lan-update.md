---
description: Update a LAN
---

# LanUpdate

## Usage

```text
ionosctl lan update [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified LAN. You can update the name, the public option for LAN and the Pcc Id to connect the LAN to a Private Cross-Connect.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId State] (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for update
  -i, --lan-id string          The unique LAN Id (required)
  -n, --name string            The name of the LAN
  -o, --output string          Desired output format [text|json] (default "text")
      --pcc-id string          The unique Id of the Private Cross-Connect the LAN will connect to
      --public                 Public option for LAN
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for LAN update to be executed
```

## Examples

```text
ionosctl lan update --datacenter-id DATACENTER_ID --lan-id LAN_ID --name NAME --public=true
```

