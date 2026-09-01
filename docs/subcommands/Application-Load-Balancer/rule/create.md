---
description: "Create an Application Load Balancer Forwarding Rule"
---

# ApplicationloadbalancerRuleCreate

## Usage

```text
ionosctl compute applicationloadbalancer rule create [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `rule` command:

```text
[r forwardingrule]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a forwarding rule (a listener) on a specified Application Load Balancer. The rule binds a protocol, an inbound IP and a port that the balancer will accept client connections on.

The --listener-ip must be one of the ALB's own IPs on its listener LAN. For an HTTPS listener, additionally pass --server-certificates so the balancer can present a certificate during the TLS handshake. After the rule exists, attach HTTP rules to it (`alb rule httprule add`) to define how matching requests are routed.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Listener Ip
* Listener Port

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --client-timeout int                  The maximum time in milliseconds to wait for the client to acknowledge or send data before the connection is closed. (default 50)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ClientTimeout ServerCertificates State]
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Level of detail for response objects (default 1)
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Maximum number of items to return per request (default 50)
      --listener-ip ip                      The inbound IP the balancer listens on. Must be one of the ALB's own --ips assigned on its listener LAN. (required)
      --listener-port int                   The inbound TCP port the balancer listens on; valid range is 1 to 65535 (typically 80 for HTTP or 443 for HTTPS). (required) (default 8080)
  -n, --name string                         The name of the Application Load Balancer forwarding rule. (default "Unnamed Forwarding Rule")
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -p, --protocol string                     The listener protocol. HTTP is the only supported value (the ALB is a layer-7 balancer); HTTPS listeners are configured by using HTTP together with --server-certificates. (default "HTTP")
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
      --server-certificates strings         IDs of server certificates (managed by the IONOS Certificate Manager) that the balancer presents to clients during the TLS handshake. Required to serve HTTPS on this listener.
  -t, --timeout int                         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create an HTTP listener on port 80
ionosctl compute applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --name "http-listener" --listener-ip 192.0.2.10 --listener-port 80

# Create an HTTPS listener on port 443 with a server certificate and a longer client timeout
ionosctl compute applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --name "https-listener" --listener-ip 192.0.2.10 --listener-port 443 --server-certificates CERTIFICATE_ID --client-timeout 60000 --wait
```

