---
description: Delete a Server
---

# ServerDelete

## Usage

```text
ionosctl server delete [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Server from a Virtual Data Center.

NOTE: This will not automatically remove the storage Volumes attached to a Server.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -a, --all                    Delete all Servers form a virtual Datacenter.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server deletion to be executed
```

## Examples

```text
ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID

ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --force
```

