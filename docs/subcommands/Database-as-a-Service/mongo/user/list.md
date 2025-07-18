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
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
  -M, --max-results int32     The maximum number of elements to return
      --no-headers            Don't print table headers when table output is used
      --offset int32          Skip a certain number of results
  -o, --output string         Desired output format [text|json|api-json] (default "text")
  -q, --quiet                 Quiet output
  -v, --verbose               Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo user list
ionosctl dbaas mongo user list --cluster-name <cluster-name>,
ionosctl dbaas mongo user list --cluster-id <cluster-id>
```

