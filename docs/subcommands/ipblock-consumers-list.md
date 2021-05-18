---
description: List IpConsumers
---

# IpblockConsumersList

## Usage

```text
ionosctl ipblock consumers list [flags]
```

## Description

Use this command to get a list of Resources where each IP address from an IpBlock is being used.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [Ip,NicId,ServerId,DatacenterId,K8sNodePoolId,K8sClusterId])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force               Force command to execute without user input
  -h, --help                help for list
      --ipblock-id string   The unique IpBlock Id (required)
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl ipblock consumers list --ipblock-id 564f4984-8349-40c1-bcd8-ba177ebf2fb6
```

