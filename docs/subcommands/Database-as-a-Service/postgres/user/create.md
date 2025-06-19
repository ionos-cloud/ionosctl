---
description: "Create a user"
---

# DbaasPostgresUserCreate

## Usage

```text
ionosctl dbaas postgres user create [flags]
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

## Description

Create a new user in the given cluster

## Options

```text
  -i, --cluster-id string   The ID of the Postgres cluster
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The password for the user
  -q, --quiet               Quiet output
      --user string         The name of the user
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres user create --cluster-id <cluster-id> --user <user> --password <password>
```

