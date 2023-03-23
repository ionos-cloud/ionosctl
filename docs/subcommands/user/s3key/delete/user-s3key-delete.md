---
description: Delete a S3Key
---

# UserS3keyDelete

## Usage

```text
ionosctl user s3key delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specific S3Key of an User.

Required values to run command:

* User Id
* S3Key Id

## Options

```text
  -a, --all                Delete all the S3Keys of an User.
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -i, --s3key-id string    The unique User S3Key Id (required)
  -t, --timeout int        Timeout option for Request for User S3Key deletion [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -w, --wait-for-request   Wait for Request for User S3Key deletion to be executed
```

## Examples

```text
ionosctl user s3key delete --user-id USER_ID --s3key-id S3KEY_ID --force
```

