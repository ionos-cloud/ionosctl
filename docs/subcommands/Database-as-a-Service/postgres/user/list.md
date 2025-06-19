---
description: "List users"
---

# DbaasPostgresUserList

## Usage

```text
ionosctl dbaas postgres user list [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
```

For `user` command:

```text
[usr u users]
```

For `list` command:

```text
[ls]
```

## Description

List all users in the given cluster

## Options

```text
  -i, --cluster-id string   The ID of the Postgres cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Username System ClusterId]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --system              List system users along with regular users
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres user list --cluster-id <cluster-id>
```

