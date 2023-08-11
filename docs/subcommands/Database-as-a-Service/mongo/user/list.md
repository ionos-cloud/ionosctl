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
[mongodb mdb m]
```

For `list` command:

```text
[l ls]
```

## Description

Retrieves a list of MongoDB users. You can either list users of a certain cluster (--cluster-id), or all clusters with an optional partial-match name filter (--cluster-name)

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string     
      --cluster-name string   When listing all users, you can optionally filter by partial-match cluster name
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Username CreatedBy Roles]
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
  -M, --max-results int32     The maximum number of elements to return
      --no-headers            When using text output, don't print headers
      --offset int32          Skip a certain number of results
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
  -v, --verbose               Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo user list
ionosctl dbaas mongo user list --cluster-name <cluster-name>,
ionosctl dbaas mongo user list --cluster-id <cluster-id>
```

