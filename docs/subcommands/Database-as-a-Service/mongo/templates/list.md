---
description: "List Mongo Templates"
---

# DbaasMongoTemplatesList

## Usage

```text
ionosctl dbaas mongo templates list [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `templates` command:

```text
[t]
```

For `list` command:

```text
[l ls]
```

## Description

Retrieves a list of valid templates. These templates can be used to create MongoDB clusters; they contain properties, such as number of cores, RAM, and the storage size.

## Options

```text
<<<<<<< HEAD
  -u, --api-url string   Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [TemplateId Name Edition Cores StorageSize RAM] (default [TemplateId,Name,Edition,Cores,StorageSize,RAM])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
      --query string     JMESPath query string to filter the output
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
=======
  -u, --api-url string    Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [TemplateId Name Edition Cores StorageSize RAM] (default [TemplateId,Name,Edition,Cores,StorageSize,RAM])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
>>>>>>> 8e970fd7 (remove deprecated 'D' for 'datacenter-id' only on psql)
```

## Examples

```text
ionosctl dbaas mongo templates list
```

