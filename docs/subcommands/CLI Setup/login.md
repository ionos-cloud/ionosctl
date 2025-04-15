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
Note: If using '--token', you can skip verifying the used token with '--skip-verify'

If you specify a custom '--api-url', the custom URL will also be saved to the config file when you login successfully and used for future API calls.

If you use a username and password, this command will use these credentials to generate a token that will be saved in the config file. Please keep this token safe.

To find your config file location, use 'ionosctl cfg location'. If you want to use a different config file, use the '--config' global option. Changing the permissions of the config file will cause it to no longer work.

If a config file already exists, you will be asked to replace it. You can skip this verification with '--force'

AUTHENTICATION ORDER
ionosctl uses a layered approach for authentication, prioritizing sources in this order:
  1. Global flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.

## Options

```text
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -p, --password string    Password to authenticate. Will be used to generate a token
  -q, --quiet              Quiet output
      --skip-verify        Forcefully write the provided token to the config file without verifying if it is valid. Note: --token is required
      --timeout duration   Timeout for waiting for resource to reach desired state (default 1m0s)
  -t, --token string       Token to authenticate. If used, will be saved to the config file without generating a new token. Note: mutually exclusive with --user and --password
      --user string        Username to authenticate. Will be used to generate a token
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed
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

