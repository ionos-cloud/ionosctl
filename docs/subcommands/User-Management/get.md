---
description: "Get a User"
---

# UserGet

## Usage

```text
ionosctl user get [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific User.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32      Controls the detail depth of the response objects. Max depth is 10.
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -i, --user-id string   The unique User Id (required)
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl user get --user-id USER_ID
```

