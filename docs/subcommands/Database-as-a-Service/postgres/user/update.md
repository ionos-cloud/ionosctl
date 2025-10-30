---
description: "Update a user"
---

# DbaasPostgresUserUpdate

## Usage

```text
ionosctl dbaas postgres user update [flags]
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

Update the specified user from the given cluster. Only changing their password is supported. System users cannot be patched.

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --cluster-id string   The ID of the Postgres cluster
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The password for the user
  -q, --quiet               Quiet output
      --user string         The name of the user
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres user update --cluster-id <cluster-id> --user <user> --password <password>
```

