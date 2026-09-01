---
description: "Check whether a registry name is available"
---

# ContainerRegistryNames

## Usage

```text
ionosctl container-registry names [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `names` command:

```text
[check name n]
```

## Description

Check whether a desired registry name is still free. Registry names are globally unique across all IONOS customers because they form the public hostname, so this check should be run before 'container-registry registry create'.

A valid name is 3-63 characters, lowercase letters/digits/dashes only, starting with a letter and ending with a letter or digit (regex ^[a-z][-a-z0-9]{1,61}[a-z0-9]$). The command reports whether the name is available, already taken, or invalid.

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -n, --name string       The desired registry name to check. Must be 3-63 chars, lowercase letters/digits/dashes, starting with a letter (regex ^[a-z][-a-z0-9]{1,61}[a-z0-9]$) (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
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
ionosctl container-registry name --name my-registry
```

