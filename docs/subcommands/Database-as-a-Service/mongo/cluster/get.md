---
description: "Get a Mongo Cluster by ID"
---

# DbaasMongoClusterGet

## Usage

```text
ionosctl dbaas mongo cluster get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get a Mongo Cluster by ID

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo cluster get --cluster-id <cluster-id>
```

