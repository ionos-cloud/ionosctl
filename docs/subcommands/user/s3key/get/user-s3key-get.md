---
description: Get a User S3Key
---

# UserS3keyGet

## Usage

```text
ionosctl user s3key get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified S3Key from a specified User.

Required values to run command:

* User Id
* S3Key Id

## Options

```text
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
      --no-headers        When using text output, don't print headers
  -i, --s3key-id string   The unique User S3Key Id (required)
      --user-id string    The unique User Id (required)
```

## Examples

```text
ionosctl user s3key get --user-id USER_ID --s3key-id S3KEY_ID
```

