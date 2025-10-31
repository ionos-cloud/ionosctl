---
description: "List Mongo Clusters"
---

# DbaasMongoClusterList

## Usage

```text
ionosctl dbaas mongo cluster list [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `cluster` command:

```text
[c]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string      Response filter to list only the Mongo Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo cluster list
```

