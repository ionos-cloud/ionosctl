---
description: List User S3Keys
---

# UserS3KeyList

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
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [S3KeyId Active] (default [S3KeyId,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
      --user-id string   The unique User Id (required)
```

## Examples

```text
ionosctl user s3key list --user-id USER_ID
```

