---
description: "Get a User S3Key"
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
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [S3KeyId Active] (default [S3KeyId,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --no-headers        When using text output, don't print headers
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -i, --s3key-id string   The unique User S3Key Id (required)
      --user-id string    The unique User Id (required)
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl user s3key get --user-id USER_ID --s3key-id S3KEY_ID
```

