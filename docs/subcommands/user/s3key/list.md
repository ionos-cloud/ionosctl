---
description: List User S3Keys
---

# UserS3keyList

## Usage

```text
ionosctl user s3key list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of S3Keys of a specified User.

Required values to run command:

* User Id

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --user-id string      The unique User Id (required)
```

## Examples

```text
ionosctl user s3key list --user-id USER_ID
```

