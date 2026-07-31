---
description: "Add an access scope to a token"
---

# ContainerRegistryTokenScopeAdd

## Usage

```text
ionosctl container-registry token scope add [flags]
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

For `add` command:

```text
[a ad]
```

## Description

Add a scope to an existing token, granting it a set of actions on a target resource. The scope is appended to the token's existing scopes (the token is patched, not replaced), so its password and other scopes are preserved.

--name is the target resource (a repository name/path, or '*' for all), --type is the kind of target ('repository', 'namespace' or 'registry'), and --actions is the comma-separated list of allowed operations ('pull', 'push', 'delete', or '*').

## Options

```text
  -a, --actions strings      Comma-separated operations the token may perform on the target, e.g. pull, push, delete (or '*' for all) (required)
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId DisplayName Type Actions]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -n, --name string          Target resource of the scope: a repository name/path, or '*' for all repositories (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   The unique ID of the registry that owns the token (required)
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
      --token-id string      The unique ID of the token to add the scope to
  -y, --type string          Kind of target the --name refers to: 'repository' (one repo), 'namespace' (a repo path prefix) or 'registry' (the whole registry) (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Grant pull+push on a single repository
ionosctl container-registry token scope add --registry-id REGISTRY_ID --token-id TOKEN_ID --name my-app --type repository --actions pull,push

# Grant read-only (pull) access to every repository in the registry
ionosctl container-registry token scope add --registry-id REGISTRY_ID --token-id TOKEN_ID --name "*" --type repository --actions pull
```

