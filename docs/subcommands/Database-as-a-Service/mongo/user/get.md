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
[m mdb mongodb mg]
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
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -d, --database string     The authentication database
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --user string         The authentication username
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo user get
```

