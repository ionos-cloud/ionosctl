---
description: "Delete a Target Group"
---

# TargetgroupDelete

## Usage

```text
ionosctl compute targetgroup delete [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `delete` command:

```text
[d]
```

## Description

Delete a Target Group. This removes the backend pool definition itself. Deleting a group that is still referenced by an ALB forwarding rule may be rejected or leave that rule pointing at a missing group, so detach it from any FORWARD rules first. Use --all to delete every Target Group in the contract.

Required values to run command:

* Target Group Id

## Options

```text
  -a, --all                     Delete all Target Groups in the contract instead of a single one. Cannot be combined with --targetgroup-id.
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval State Retries Path Method MatchType Response Regex Negate]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Maximum number of items to return per request (default 50)
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Delete one target group
ionosctl compute targetgroup delete --targetgroup-id TARGET_GROUP_ID --force

# Delete every target group in the contract
ionosctl compute targetgroup delete --all --force
```

