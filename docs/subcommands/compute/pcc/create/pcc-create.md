---
description: Create a Private Cross-Connect
---

# PccCreate

## Usage

```text
ionosctl pcc create [flags]
```

## Aliases

For `create` command:

```text
[c]
```

## Description

Use this command to create a Private Cross-Connect. You can specify the name and the description for the Private Cross-Connect.

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string   The description for the Private Cross-Connect
  -n, --name string          The name for the Private Cross-Connect (default "Unnamed PrivateCrossConnect")
  -t, --timeout int          Timeout option for Request for Private Cross-Connect creation [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Private Cross-Connect creation to be executed
```

## Examples

```text
ionosctl pcc create --name NAME --description DESCRIPTION --wait-for-request
```

