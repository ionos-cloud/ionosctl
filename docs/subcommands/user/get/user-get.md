---
description: Get a User
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
  -D, --depth int32      Controls the detail depth of the response objects. Max depth is 10.
      --no-headers       When using text output, don't print headers
  -i, --user-id string   The unique User Id (required)
```

## Examples

```text
ionosctl user get --user-id USER_ID
```

