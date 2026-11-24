---
description: "List MongoDB Instance Templates"
---

# DbaasMongoV2TemplateList

## Usage

```text
ionosctl dbaas mongo-v2 template list [flags]
```

## Aliases

For `mongo-v2` command:

```text
[mongodb-v2 mg-v2]
```

For `template` command:

```text
[templates t]
```

For `list` command:

```text
[ls l]
```

## Description

Use this command to retrieve a list of MongoDB Instance Templates. You can filter by edition using `--edition`.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mongodbv2' and env var 'IONOS_API_URL' (default "https://mongodb.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [TemplateId Name Edition Cores Ram StorageSize]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
      --edition string    Response filter to list only the templates of this edition. Can be one of: ENTERPRISE, BUSINESS
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int32       The maximum number of elements to return (default 100)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: gb/bhx, us/ewr, gb/lhr, fr/par, es/vit, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --no-headers        Don't print table headers when table output is used
      --offset int32      The first element to return
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo-v2 template list --edition BUSINESS
```

