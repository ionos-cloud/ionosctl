---
description: Add a Label on a Server
---

# AddLabel

## Usage

```text
ionosctl server add-label [flags]
```

## Description

Use this command to add/create a Label on Server. You must specify the key and the value for the Label.

Required values to run command:

* Data Center Id
* Server Id
* Label Key
* Label Value

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for add-label
      --ignore-stdin           Force command to execute without user input
      --label-key string       The unique Label Key [Required flag]
      --label-value string     The unique Label Value [Required flag]
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
```

## Examples

```text
ionosctl server add-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --server-id 27dde318-f0d4-4f97-a04d-9dafe4a89637 --label-key test --label-value test
Key    Value
test   test
```

