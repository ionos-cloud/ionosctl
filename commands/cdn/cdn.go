package cdn

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/distribution"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "cdn",
			Short: "Manage CDN distributions",
			Long: `The Content Delivery Network (CDN) service caches your content at IONOS edge locations (Points of Presence) around the world and serves it to visitors from the nearest edge, reducing latency and offloading traffic from your origin servers. It also provides Layer 7 DDoS protection and an optional Web Application Firewall (WAF).

The single resource you manage here is a DISTRIBUTION. A distribution ties together three things:
  - a DOMAIN (the public hostname visitors use, e.g. cdn.example.com),
  - an optional CERTIFICATE binding (a certificate manager UUID used to terminate HTTPS for that domain), and
  - a list of ROUTING RULES that match requests by URL path prefix and forward each to an upstream origin, deciding per-rule whether to cache, apply the WAF, rate-limit, and geo-restrict.

Once a distribution is AVAILABLE, point your domain's DNS at the CDN (typically a CNAME to the value shown in the distribution's metadata) so that visitor traffic flows through the edge network.`,
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(distribution.Command())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.CDN}, constants.CDNApiRegionalURL, constants.CDNLocations)
}
