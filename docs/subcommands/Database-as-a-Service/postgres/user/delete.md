---
description: "Delete user"
---

# DbaasPostgresUserDelete

## Usage

```text
ionosctl dbaas postgres user delete [flags]
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

For `delete` command:

```text
[del]
```

## Description

Delete the specified user from the given cluster

## Options

```text
  -i, --cluster-id string   The ID of the Postgres cluster
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --user string         The name of the user
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres user delete --cluster-id <cluster-id> --user <user>
```

