---
description: "Create or replace the CORS configuration for a bucket"
---

# ObjectStorageCorsPut

## Usage

```text
ionosctl object-storage cors put [flags]
```

## Aliases

For `object-storage` command:

```text
[os]
```

For `put` command:

```text
[p]
```

## Description

Create or replace the CORS configuration for a bucket. The configuration must be provided as a path to a JSON file via --json-properties. Use --json-properties-example to see an example CORS configuration.

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [AllowedOrigins AllowedMethods AllowedHeaders ExposeHeaders MaxAgeSeconds ID]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file containing the CORS configuration
      --json-properties-example   Print an example CORS configuration JSON and exit
      --limit int                 Maximum number of items to return per request (default 50)
  -n, --name string               Name of the bucket (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl object-storage cors put --name my-bucket --json-properties cors.json
ionosctl object-storage cors put --json-properties-example
```

