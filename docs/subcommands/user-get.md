---
description: Get a User
---

# UserGet

## Usage

```text
ionosctl user get [flags]
```

## Aliases

For `user` command:
```text
[u]
```

For `get` command:
```text
[g]
```

## Description

Use this command to retrieve details about a specific User.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for get
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -i, --user-id string   The unique User Id (required)
```

## Examples

```text
ionosctl user get --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
2470f439-1d73-42f8-90a9-f78cf2776c74   test1       test1      testrandom12@ionos.com   false           false          false           a74101e7c1948450432d5b6512f2712c   true
```

