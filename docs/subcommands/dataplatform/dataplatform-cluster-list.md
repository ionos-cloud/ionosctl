---
description: List Data Platform Clusters
---

# DataplatformClusterList

## Usage

```text
ionosctl dataplatform cluster list [flags]
```

## Aliases

For `dataplatform` command:

```text
[dp]
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

Use this command to retrieve a list of Data Platform Clusters provisioned under your account. You can filter the result based on Data Platform Name using `--name` option.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [ClusterId,Name,DataPlatformVersion,MaintenanceWindow,DatacenterId,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --name string      Response filter to list only the Data Platform Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform cluster list
```

