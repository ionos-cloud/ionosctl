---
description: "List User S3Keys"
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
  -u, --api-url string      Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [S3KeyId Active SecretKey] (default [S3KeyId,Active,SecretKey])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --user-id string      The unique User Id (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl user s3key list --user-id USER_ID
```

