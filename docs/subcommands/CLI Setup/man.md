---
description: "Generate manpages for ionosctl"
---

# Man

## Usage

```text
ionosctl man [flags]
```

## Aliases

For `man` command:

```text
[manpages]
```

## Description

WARNING: This command is only supported on Linux.

The 'man' command allows you to generate manpages for ionosctl in a given directory. By default, the manpages will be compressed using gzip, but you can skip this step by using the '--skip-compression' flag.
In order to install the manpages, there are a few steps you need to follow:
- Decide where you would like to install the manpages. You can check which directories are available to you by running 'manpath'. If you want to install the manpages to a directory that is not listed, you can add a new entry to '~/.manpath' (see 'man 5 manpath' for how to do it). The directory must contain subdirectories for each section (e.g. 'man1', 'man5', etc.).
- Copy the manpages to the installation directory, in the 'man1' section.
- Run 'sudo mandb' to update the 'man' internal database.

After following these steps, you should be able to use 'man ionosctl' to access the manpages.

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
      --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --skip-compression    Skip compressing manpages with gzip, just generate them
      --target-dir string   Target directory where manpages will be generated. Must be an absolute path (default "/tmp/ionosctl-man")
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl man --target-dir /tmp/ionosctl-man
```

