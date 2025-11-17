---
description: "Retrieves a list of MongoDB users. You can either list users of a certain cluster (--cluster-id), or all clusters with an optional partial-match name filter (--cluster-name)"
---

# DbaasMongoUserList

## Usage

```text
ionosctl dbaas mongo user list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Retrieves a list of MongoDB users. You can either list users of a certain cluster (--cluster-id), or all clusters with an optional partial-match name filter (--cluster-name)

## Options

```text
  -u, --api-url string        Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string     
      --cluster-name string   When listing all users, you can optionally filter by partial-match cluster name
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Username CreatedBy Roles]
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int             Level of detail for response objects (default 1)
      --filters strings       Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --limit int             Maximum number of items to return per request (default 50)
      --no-headers            Don't print table headers when table output is used
      --offset int            Number of items to skip before starting to collect the results
      --order-by string       Property to order the results by
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo user list
ionosctl dbaas mongo user list --cluster-name <cluster-name>,
ionosctl dbaas mongo user list --cluster-id <cluster-id>
```

