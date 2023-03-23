---
description: Create a S3Key for a User
---

# UserS3keyCreate

## Usage

```text
ionosctl user s3key create [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `s3key` command:

```text
[k s3k]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a S3Key for a particular User.

Note: A maximum of five S3 keys may be created for any given user.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* User Id

## Options

```text
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -t, --timeout int        Timeout option for Request for User S3Key creation [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -w, --wait-for-request   Wait for the Request for User S3Key creation to be executed
```

## Examples

```text
ionosctl user s3key create --user-id USER_ID
```

