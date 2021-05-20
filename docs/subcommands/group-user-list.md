---
description: List Users from a Group
---

# GroupUserList

## Usage

```text
ionosctl group user list [flags]
```

## Aliases

For `group` command:
```text
[g]
```

For `user` command:
```text
[u]
```

## Description

Use this command to get a list of Users from a specific Group.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              help for list
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl group user list --group-id 45ba215b-6897-40b6-879c-cbadb527cefd 
UserId                                 Firstname   Lastname   Email                    S3CanonicalUserId                  Administrator   ForceSecAuth   SecAuthActive   Active
62599641-aa2d-4ecc-bdc4-118f5f39f23d   test        test       testrandom53@gmail.com   f670112b3e74038b51db78d5836d7854   false           false          false           true
```

