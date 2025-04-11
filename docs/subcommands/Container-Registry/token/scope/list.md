<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> a3be15e5 (fix: timeout add -t only for commands where -t not exist)
---
description: "Get a token scopes"
---

# ContainerRegistryTokenScopeList

## Usage

```text
ionosctl container-registry token scope list [flags]
```

## Aliases

For `token` command:

```text
[t tokens]
```

For `scope` command:

```text
[s scopes]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list all scopes of a token of a container registry.

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
<<<<<<< HEAD
                             Available columns: [ScopeId DisplayName Type Actions]
=======
                             Available columns: [ScopeId TokenId DisplayName Type Actions]
>>>>>>> a3be15e5 (fix: timeout add -t only for commands where -t not exist)
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
<<<<<<< HEAD
  -t, --token-id string      Token ID
  -v, --verbose              Print step-by-step process when running command
=======
      --timeout duration     Timeout for waiting for resource to reach desired state (default 1m0s)
  -t, --token-id string      Token ID
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed
>>>>>>> a3be15e5 (fix: timeout add -t only for commands where -t not exist)
```

## Examples

```text
ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

<<<<<<< HEAD
=======
>>>>>>> c6629594 (doc: regen docs)
=======
>>>>>>> a3be15e5 (fix: timeout add -t only for commands where -t not exist)
