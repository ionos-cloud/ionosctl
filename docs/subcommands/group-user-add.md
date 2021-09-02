---
description: Add User to a Group
---

# GroupUserAdd

## Usage

```text
ionosctl group user add [flags]
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

For `add` command:

```text
[a]
```

## Description

Use this command to add an existing User to a specific Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              Print usage
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -i, --user-id string    The unique User Id (required)
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl group user add --group-id GROUP_ID --user-id USER_ID
```

