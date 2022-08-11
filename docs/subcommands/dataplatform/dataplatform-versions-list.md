---
description: List Managed Data Platform API versions
---

# DataplatformVersionsList

## Usage

```text
ionosctl dataplatform versions list [flags]
```

## Aliases

For `dataplatform` command:

```text
[dp]
```

For `versions` command:

```text
[v]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Managed Data Platform API versions.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [DataPlatformVersion] (default [DataPlatformVersion])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform versions list --cluster-id CLUSTER_ID
```

