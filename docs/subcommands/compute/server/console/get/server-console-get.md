---
description: Get the Remote Console URL to access a Server
---

# ServerConsoleGet

## Usage

```text
ionosctl server console get [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `console` command:

```text
[url]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get the Server Remote Console link.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --no-headers             When using text output, don't print headers
  -i, --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl server console get --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

