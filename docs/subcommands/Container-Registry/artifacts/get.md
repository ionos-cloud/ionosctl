---
description: "Retrieve an artifacts"
---

# ContainerRegistryArtifactsGet

## Usage

```text
ionosctl container-registry artifacts get [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `artifacts` command:

```text
[a art artifact]
```

## Description

Retrieve an artifact from a repository

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --artifact-id string   ID/digest of the artifact
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Repository PushCount PullCount LastPushed TotalVulnerabilities FixableVulnerabilities MediaType URN RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
      --repository string    Name of the repository to retrieve artifact from
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry artifacts get
```

