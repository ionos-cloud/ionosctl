---
description: "Get a MongoDB user"
---

# DbaasMongoUserGet

## Usage

```text
ionosctl dbaas mongo user get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get a MongoDB user

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --database string     The authentication database
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --user string         The authentication username
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo user get
```

