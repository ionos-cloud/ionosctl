---
description: ionosctl is a command line interface (CLI) for the Ionos Cloud
---

# Ionosctl

## Usage

```text
ionosctl [command]
```

## Description


        _                                         __     __
       (_)  ____    ____   ____    _____  _____  / /_   / /
      / /  / __ \  / __ \ / __ \  / ___/ / ___/ / __/  / /
     / /  / /_/ / / / / // /_/ / (__  ) / /__  / /_   / /
    /_/   \____/ /_/ /_/ \____/ /____/  \___/  \__/  /_/

The IonosCTL wraps the Ionos Cloud API allowing you to interact with it from a command-line interface.
The command `ionosctl` is the root command that all other commands are attached to.
IonosCTL supports json format for all output commands by setting `--output=json` option.

Note: if error, it returns exit code 1.


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

