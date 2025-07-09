---
description: "Remove User from a Group"
---

# GroupUserRemove

## Usage

```text
ionosctl group user remove [flags]
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

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a User from a Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -a, --all               Remove all Users from a group.
  -u, --api-url string    Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              Print usage
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -i, --user-id string    The unique User Id (required)
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl group user remove --group-id GROUP_ID --user-id USER_ID
```

