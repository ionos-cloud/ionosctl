---
description: Delete a Private Cross-Connect
---

# PccDelete

## Usage

```text
ionosctl pcc delete [flags]
```

## Aliases

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -a, --all                Delete all Private Cross-Connects.
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -i, --pcc-id string      The unique Private Cross-Connect Id (required)
  -t, --timeout int        Timeout option for Request for Private Cross-Connect deletion [seconds] (default 60)
  -w, --wait-for-request   Wait for the Request for Private Cross-Connect deletion to be executed
```

## Examples

```text
ionosctl pcc delete --pcc-id PCC_ID --wait-for-request
```

