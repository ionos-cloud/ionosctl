---
description: "List Dataplatform Clusters"
---

# DataplatformClusterList

## Usage

```text
ionosctl dataplatform cluster list [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
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

List Dataplatform Clusters

## Options

```text
      --cols strings    Set of columns to be printed on output 
                        Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -c, --config string   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force           Force command to execute without user input
  -h, --help            Print usage
  -n, --name string     Response filter to list only the clusters which include the specified name. case insensitive
      --no-headers      Don't print table headers when table output is used
  -o, --output string   Desired output format [text|json|api-json] (default "text")
  -q, --quiet           Quiet output
  -v, --verbose         Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform cluster list
```

