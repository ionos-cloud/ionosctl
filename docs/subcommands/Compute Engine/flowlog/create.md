---
description: "Create a FlowLog on a NIC"
---

# FlowlogCreate

## Usage

```text
ionosctl compute flowlog create [flags]
```

## Aliases

For `flowlog` command:

```text
[fl]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Flow Log on a specific NIC. The Flow Log streams traffic metadata for that NIC into an existing IONOS Object Storage bucket.

The bucket named by --s3bucket must already exist and be in a location that supports Flow Log delivery; this command does not create it. Choose which traffic to capture with --action (ACCEPTED / REJECTED / ALL) and --direction (INGRESS / EGRESS / BIDIRECTIONAL).

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

NOTE: Disable/delete the Flow Log before deleting the Object Storage bucket it writes to.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Target S3 Bucket Name

## Options

```text
  -a, --action string          Which traffic to capture by firewall outcome: ACCEPTED (allowed traffic), REJECTED (blocked traffic), or ALL (both) (default "ALL")
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique ID of the Virtual Data Center that holds the server and NIC (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --direction string       Which traffic direction to capture: INGRESS (inbound to the NIC), EGRESS (outbound from the NIC), or BIDIRECTIONAL (both) (default "BIDIRECTIONAL")
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            A display name for the Flow Log (default "Unnamed FlowLog")
      --nic-id string          The unique ID of the NIC whose traffic will be captured (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -b, --s3bucket string        Name of an EXISTING IONOS Object Storage (S3) bucket that will receive the Flow Log files. The bucket must already exist (required)
      --server-id string       The unique ID of the server that owns the NIC (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Capture all traffic in both directions
ionosctl compute flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name "audit-all" --s3bucket my-flowlog-bucket

# Capture only rejected inbound traffic (useful for security review) and wait until ready
ionosctl compute flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name "blocked-ingress" --action REJECTED --direction INGRESS --s3bucket my-flowlog-bucket --wait
```

