---
description: "Returns latest Dataplatform Cluster version, which you can use in cluster creation."
---

# DataplatformVersionActive

## Usage

```text
ionosctl dataplatform version active [flags]
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

For `active` command:

```text
[latest last]
```

## Description

Returns latest Dataplatform Cluster version, which you can use in cluster creation.

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'dataplatform' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

