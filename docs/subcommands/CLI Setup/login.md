---
description: "Use credentials to generate a config file in `ionosctl cfg location`, or use '--example' to generate a sample endpoints YAML config"
---

# ConfigLogin

## Usage

```text
ionosctl config login [flags]
```

## Aliases

For `config` command:

```text
[cfg]
```

## Description

Generate a YAML file aggregating all product endpoint information at 'ionosctl cfg location' using the public OpenAPI index.

If using '--example', this command prints the config to stdout without any authentication step.

You can filter by version (--filter-version), whitelist (--whitelist) or blacklist (--blacklist) specific APIs,
and customize the names of the APIs in the config file using --custom-names.

There are three ways you can authenticate with the IONOS Cloud APIs:
  1. Interactive mode: Prompts for username and password, and generates a token that will be saved in the config file.
  2. Use the '--user' and '--password' flags: Used to generate a token that will be saved in the config file.
  3. Use the '--token' flag: Provide an authentication token.
Notes:
  - If using '--example', the authentication step is skipped


## Options

```text
  -u, --api-url string                Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
      --blacklist strings             Comma-separated list of API names to exclude (default [object-storage-user-owned-buckets,object-storage-contract-owned-buckets,identity-federation,identity-provider,identity-policy,inference-modelhub,inference-openai,quota,reseller,tagging])
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --custom-names stringToString   Define custom names for each spec (default <Overriden with sdk-go-bundle product names: [authentication=auth, certificatemanager=cert, cloud=compute, object‑storage=objectstorage, object‑storage‑management=objectstoragemanagement, mongodb=mongo, postgresql=psql]>)
      --environment string            Environment to use (default "prod")
      --example                       Print an example YAML config file to stdout and skip authentication step
      --filter-version string         Filter by major spec version (e.g. v1)
  -f, --force                         Force command to execute without user input
  -h, --help                          Print usage
      --no-headers                    Don't print table headers when table output is used
  -o, --output string                 Desired output format [text|json|api-json] (default "text")
  -p, --password string               Password to authenticate with. Will be used to generate a token
      --profile-name string           Name of the profile to use (default "user")
  -q, --quiet                         Quiet output
      --skip-verify                   Forcefully write the provided token to the config file without verifying if it is valid. Note: --token is required
  -t, --token string                  Token to authenticate with. If used, will be saved directly to the config file. Note: mutually exclusive with --user and --password
      --user string                   Username to authenticate with. Will be used to generate a token
  -v, --verbose                       Print step-by-step process when running command
      --version float                 Version of the config file to use (default 1)
      --whitelist strings             Comma-separated list of API names to include
```

## Examples

```text

# Print an example YAML configuration file to stdout
ionosctl config login --example

# Login interactively, and generate a YAML config file with filters, to 'ionosctl config location'
ionosctl endpoints generate --filter-version=v1 \
  --whitelist=vpn,psql --blacklist=billing

# Specify a token, a config version, a custom profile name, and a custom environment
ionosctl config login --token $IONOS_TOKEN \
  --version=1.1 --profile-name=my-custom-profile --environment=dev

```

