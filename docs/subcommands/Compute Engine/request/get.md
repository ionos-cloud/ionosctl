---
description: "Get a Request"
---

# RequestGet

## Usage

```text
ionosctl request get [flags]
```

## Aliases

For `request` command:

```text
[req]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Request.

Required values to run command:

* Request Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [RequestId CreatedDate CreatedBy Method Status Message Url Body Targets] (default [RequestId,CreatedDate,Method,Status,Message,Targets])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -i, --request-id string   The unique Request Id (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl request get --request-id REQUEST_ID
```

