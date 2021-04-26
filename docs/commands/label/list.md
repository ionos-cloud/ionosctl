---
description: List Labels from all Resources
---

# List

## Usage

```text
ionosctl label list [flags]
```

## Description

Use this command to list all Labels from all Resources under your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl label list 
Key       Value            ResourceType   ResourceId
test      testserver       server         27dde318-f0d4-4f97-a04d-9dafe4a89637
test      testdatacenter   datacenter     ed612a0a-9506-4b56-8d1b-ce2b04090f19
test      testsnapshot     snapshot       df7f4ad9-b942-4e79-939d-d1c10fb6fbff
```

