---
description: "Delete a MongoDB user"
---

# DbaasMongoUserDelete

## Usage

```text
ionosctl dbaas mongo user delete [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `user` command:

```text
[u]
```

For `delete` command:

```text
[g]
```

## Description

Delete a MongoDB user

## Options

```text
  -a, --all                 Delete all users in a cluster
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -d, --database string     The authentication database
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --name string         The authentication username
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo user delete
```

