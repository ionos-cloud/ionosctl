---
description: Get a Label
---

# Get

## Usage

```text
ionosctl label get [flags]
```

## Description

Use this command to get information about a specified Label using its URN. A URN is used for uniqueness of a Label and composed using `urn:label:<resource_type>:<resource_uuid>:<key>`.

Required values to run command:

* Label URN

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force              Force command to execute without user input
  -h, --help               help for get
      --label-urn string   URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
```

## Examples

```text
ionosctl label get --label-urn "urn:label:server:27dde318-f0d4-4f97-a04d-9dafe4a89637:test"
Key    Value        ResourceType   ResourceId
test   testserver   server         27dde318-f0d4-4f97-a04d-9dafe4a89637
```

