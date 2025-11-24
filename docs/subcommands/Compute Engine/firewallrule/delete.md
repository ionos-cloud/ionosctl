---
description: "Delete a FirewallRule"
---

# FirewallruleDelete

## Usage

```text
ionosctl firewallrule delete [flags]
```

## Aliases

For `firewallrule` command:

```text
[f fr firewall]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Firewall Rule from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id

## Options

```text
  -a, --all                      Delete all the Firewalls.
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [FirewallRuleId Name Protocol SourceMac SourceIP DestinationIP PortRangeStart PortRangeEnd IcmpCode IcmpType Direction IPVersion State] (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,Direction,IPVersion,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string     The unique Data Center Id (required)
      --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -i, --firewallrule-id string   The unique FirewallRule Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Pagination limit: Maximum number of items to return per request (default 50)
      --nic-id string            The unique NIC Id (required)
      --no-headers               Don't print table headers when table output is used
      --offset int               Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id (required)
  -t, --timeout int              Timeout option for Request for Firewall Rule deletion [seconds] (default 60)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request         Wait for Request for Firewall Rule deletion to be executed
```

## Examples

```text
ionosctl firewallrule delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID
```

