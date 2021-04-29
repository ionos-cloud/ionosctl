---
description: Add a Label on a IpBlock
---

# AddLabel

## Usage

```text
ionosctl ipblock add-label [flags]
```

## Description

Use this command to add/create a Label on IpBlock. You must specify the key and the value for the Label.

Required values to run command: 

* IpBlock Id 
* Label Key
* Label Value

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                Force command to execute without user input
  -h, --help                 help for add-label
      --ipblock-id string    The unique IpBlock Id (required)
      --label-key string     The unique Label Key (required)
      --label-value string   The unique Label Value (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl ipblock add-label --ipblock-id 379a995b-f285-493e-a56a-f32e1cb6dd06 --label-key test --label-value testipblock
Key    Value
test   testipblock
```

