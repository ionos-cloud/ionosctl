---
description: List Users
---

# UserList

## Usage

```text
ionosctl user list [flags]
```

## Aliases

For `user` command:
```text
[u]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of existing Users available on your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl user list 
UserId                                 Firstname   Lastname   Email                      Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
2470f439-1d73-42f8-90a9-f78cf2776c74   test1       test1      testrandom12@ionos.com     false           false          false           a74101e7c1948450432d5b6512f2712c   true
53d68de9-931a-4b61-b532-82f7b27afef3   test1       test1      testrandom13@ionos.com     false           false          false           8b9dd6f39e613adb7a837127edb67d38   true
```

