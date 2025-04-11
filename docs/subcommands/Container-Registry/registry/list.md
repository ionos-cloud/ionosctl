---
description: "List all Registries"
---

# ContainerRegistryRegistryList

## Usage

```text
ionosctl container-registry registry list [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `registry` command:

```text
[reg registries r]
```

For `list` command:

```text
[l ls]
```

## Description

List all managed container registries for your account

## Options

```text
<<<<<<< HEAD
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime State]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --name string      Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -t, --timeout int      Timeout in seconds for polling the request (default 60)
  -v, --verbose          Print step-by-step process when running command
  -w, --wait             Polls the request continuously until the operation is completed
=======
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [RegistryId DisplayName Location Hostname VulnerabilityScanning GarbageCollectionDays GarbageCollectionTime]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -n, --name string        Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout duration   Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed
>>>>>>> a3be15e5 (fix: timeout add -t only for commands where -t not exist)
```

## Examples

```text
ionosctl container-registry registry list
```

