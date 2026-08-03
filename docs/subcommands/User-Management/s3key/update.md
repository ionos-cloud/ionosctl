---
description: "Enable or disable a User's S3 key"
---

# UserS3keyUpdate

## Usage

```text
ionosctl compute user s3key update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Enable or disable an existing S3 key of a User by setting --s3key-active. Disabling a key immediately stops it from authenticating against Object Storage without deleting it, which makes this the safe way to rotate or temporarily suspend credentials (re-enable it later, or delete it once a replacement is in use).

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* User Id
* S3Key Id
* S3Key Active

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [S3KeyId Active SecretKey]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --s3key-active      Whether the key is active: true enables it for Object-Storage authentication, false disables it (without deleting it). E.g.: --s3key-active=true, --s3key-active=false
  -i, --s3key-id string   The unique User S3Key Id (required)
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
      --user-id string    The unique User Id (required)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false
```

