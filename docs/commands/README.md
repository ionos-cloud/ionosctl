---
description: ionosctl is a command-line interface for the Ionos Cloud
---

# Ionosctl

## Usage

```text
ionosctl [command]
```

## Description

IonosCTL is a command-line interface (CLI) for the Ionos Cloud API. 
Its main purpose is to help you manage your Ionos Cloud resources directly from your terminal.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
  -h, --help             help for ionosctl
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output
```

## See also

* [ionosctl completion](completion/)
* [ionosctl datacenter](datacenter/)
* [ionosctl lan](lan/)
* [ionosctl loadbalancer](loadbalancer/)
* [ionosctl location](location/)
* [ionosctl login](login.md)
* [ionosctl nic](nic/)
* [ionosctl request](request/)
* [ionosctl server](server/)
* [ionosctl version](version.md)
* [ionosctl volume](volume/)

