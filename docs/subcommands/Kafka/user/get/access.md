---
description: "Get a user's credentials"
---

# KafkaUserGetAccess

## Usage

```text
ionosctl kafka user get-access [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `get-access` command:

```text
[g get access]
```

## Description

Get a Kafka user's credentials including certificate, private key, and CA certificate.
By default, the command writes three PEM files to the specified output directory (or current directory if not specified):
 - <username>-cert.pem
 - <username>-key.pem
 - <username>-ca.pem

You can also use '--output json' to print the full JSON response from the API to stdout instead of writing files.

IMPORTANT: Keep these credentials secure. The private key should never be shared or exposed publicly.

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
      --cluster-id string   The ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --output-dir string   Directory to save the user's credential PEM files (default ".")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --user-id string      The ID of the user (required)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl kafka user get-access --location LOCATION --cluster-id CLUSTER_ID --user-id USER_ID 
```

