---
description: "List all Dataplatform Cluster versions, including deprecated ones. To view the latest version, use the 'version active' command"
---

# DataplatformVersionList

## Usage

```text
ionosctl dataplatform version list [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
```

For `version` command:

```text
[v]
```

For `list` command:

```text
[ls]
```

## Description

List all Dataplatform Cluster versions, including deprecated ones. To view the latest version, use the 'version active' command

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'dataplatform' and env var 'IONOS_API_URL' (default "https://api.ionos.com/dataplatform")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

