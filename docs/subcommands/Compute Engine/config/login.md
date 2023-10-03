---
description: "Use credentials to generate a config file in `ionosctl cfg location`"
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

For `login` command:

```text
[log auth]
```

## Description

The 'login' command allows you to authenticate with the IONOS Cloud APIs. There are three ways you can use it:
  1. Interactive mode: Just type 'ionosctl login' and you'll be prompted to enter your username and password.
  2. Use the '--user' and '--password' flags: Enter your credentials in the command.
  3. Use the '--token' flag: Provide an authentication token.
Note: If using '--token', you can skip verifying the used token and force save it by using '--force'

If you specify a custom '--api-url', the custom URL will also be saved to the config file when you login successfully and used for future API calls.

If you use a username and password, this command will use these credentials to generate a token that will be saved in the config file. Please keep this token safe.

To find your config file location, use 'ionosctl cfg location'. If you want to use a different config file, use the '--config' global option. Changing the permissions of the config file will cause it to no longer work.

AUTHENTICATION ORDER
ionosctl uses a layered approach for authentication, prioritizing sources in this order:
  1. Global flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Forcefully write the provided token to the config file without verifying if it is valid. Note: --token is required
  -h, --help              Print usage
  -o, --output string     Desired output format [text|json] (default "text")
  -p, --password string   Password to authenticate
  -q, --quiet             Quiet output
  -t, --token string      Token to authenticate
      --user string       Username to authenticate
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl login --user $IONOS_USERNAME --password $IONOS_PASSWORD

ionosctl login --token $IONOS_TOKEN

ionosctl login
Enter your username:
USERNAME
Enter your password:
```

