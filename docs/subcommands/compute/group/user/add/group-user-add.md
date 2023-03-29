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
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
      --group-id string   The unique Group Id (required)
  -i, --user-id string    The unique User Id (required)
```

## Examples

```text
ionosctl group user add --group-id GROUP_ID --user-id USER_ID
```

