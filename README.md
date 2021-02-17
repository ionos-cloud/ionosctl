# ionosctl

![CI](https://github.com/ionos-cloud/ionosctl/workflows/CI/badge.svg)

## Overview

IonosCTL is a tool to help you manage your Ionos Cloud resources directly from your terminal.

Run `ionosctl`, `ionosctl -h`, `ionosctl --help` or `ionosctl help`:
```text
IonosCTL is a command-line interface (CLI) for the Ionos Cloud API. 

USAGE: 
  ionosctl [command]

AVAILABLE COMMANDS:
  completion   Generate code to enable auto-completion with TAB key
  datacenter   Data Center Operations
  help         Help about any command
  lan          LAN Operations
  loadbalancer Load Balancer Operations
  location     Location Operations
  login        Authentication command for SDK
  nic          Network Interfaces Operations
  request      Request Operations
  server       Server Operations
  version      Show the current version
  volume       Volume Operations

FLAGS:
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "/home/ana/.config/ionosctl-config.json")
  -h, --help             help for ionosctl
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output

Use "ionosctl [command] --help" for more information about a command.
```

## Getting Started

Before you begin you will need to have signed-up for a Ionos Cloud account. The credentials you establish during sign-up will be used to authenticate against the Ionos Cloud API.

### Installation

## FAQ

1. How can I open a bug/feature request?

Bugs & feature requests can be open on the repository issues: [https://github.com/ionos-cloud/ionosctl/issues/new/choose](https://github.com/ionos-cloud/ionosctl/issues/new/choose)
