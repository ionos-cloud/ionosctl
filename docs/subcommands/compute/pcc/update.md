---
description: Update a Private Cross-Connect
---

# PccUpdate

## Usage

```text
ionosctl pcc update [flags]
```

## Aliases

For `update` command:

```text
[u up]
```

## Description

Use this command to update details about a specific Private Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string   The description for the Private Cross-Connect
  -n, --name string          The name for the Private Cross-Connect
  -i, --pcc-id string        The unique Private Cross-Connect Id (required)
  -t, --timeout int          Timeout option for Request for Private Cross-Connect update [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Private Cross-Connect update to be executed
```

## Examples

```text
ionosctl pcc update --pcc-id PCC_ID --description DESCRIPTION
```

