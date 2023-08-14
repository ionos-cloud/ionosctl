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
[mongodb mdb m]
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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --database string     The authentication database
  -f, --force               Skip y/n checks
  -h, --help                Print usage
      --name string         The authentication username
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo user delete
```

