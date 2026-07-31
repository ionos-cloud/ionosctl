package distribution

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "cdn",
		Resource:  "distribution",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a CDN distribution for a domain, with routing rules mapping URL prefixes to upstream origins",
		LongDesc: `Create a CDN distribution. A distribution serves a single DOMAIN and needs at least one ROUTING RULE describing where to fetch content that is not already cached.

Each routing rule matches requests whose path starts with a given prefix (e.g. "/api") and a scheme (http, https, or http/https), then forwards them to an upstream origin. Per rule you control caching, the WAF, a per-IP rate-limit class, geo-restrictions, and the SNI mode used when the CDN connects to the origin over TLS. Rules are supplied as JSON via --routing-rules; run 'ionosctl cdn ds create --routing-rules-example' to print a ready-to-edit template.

Provide --certificate-id (a Certificate Manager UUID) to terminate HTTPS for the domain; omit it for HTTP-only distributions.

Constraints (enforced by the API):
  - --domain must be a valid, unique hostname (2-253 chars), e.g. cdn.example.com.
  - Each distribution needs 1-25 routing rules; each rule's prefix is 1-128 chars and must start with "/".
  - Once AVAILABLE, point the domain's DNS (usually a CNAME) at the CDN so traffic reaches the edge.

Routing-rule JSON fields (per rule):
  - prefix:   URL path prefix to match, e.g. "/" or "/api".
  - scheme:   one of "http", "https", "http/https" (accept both).
  - upstream.host:           origin hostname to fetch uncached content from.
  - upstream.caching:        true/false; cache origin responses at the edge.
  - upstream.waf:            true/false; enable the Web Application Firewall.
  - upstream.rateLimitClass: per-IP request rate limit, one of R1, R5, R10, R25, R50, R100, R250, R500 (the number is the allowed requests/second per client IP; R1 is strictest, R500 most permissive).
  - upstream.sniMode:        "distribution" (origin cert must match the distribution's domain) or "origin" (origin cert must match upstream.host).
  - upstream.geoRestrictions: optionally EITHER {"allowList":[...]} (only these countries may access) OR {"blockList":[...]} (these countries are denied), using ISO 3166-1 alpha-2 codes (e.g. "DE", "US"). Use one list per rule, not both.`,
		Example: `# Create an HTTP-only distribution, passing routing rules inline as JSON
ionosctl cdn ds create --domain cdn.example.com --routing-rules '[{"prefix":"/","scheme":"http/https","upstream":{"host":"origin.example.com","caching":true,"waf":false,"rateLimitClass":"R500","sniMode":"origin"}}]'

# Print an editable routing-rules template, then create an HTTPS distribution from a rules file
ionosctl cdn ds create --routing-rules-example > rules.json
ionosctl cdn ds create --domain cdn.example.com --certificate-id 5a029f4a-72e5-11ec-90d6-0242ac120003 --routing-rules rules.json`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagCDNDistributionDomain, constants.FlagCDNDistributionRoutingRules},
				[]string{constants.FlagCDNDistributionRoutingRulesExample}); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagCDNDistributionRoutingRulesExample)) {
				fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", RoutingRuleExample)
				return nil
			}

			input := &cdn.DistributionProperties{}
			if err := setPropertiesFromFlags(c, input); err != nil {
				return err
			}

			id := uuidgen.Must()
			res, _, err := client.Must().CDNClient.DistributionsApi.DistributionsPut(context.Background(), id).
				DistributionUpdate(cdn.DistributionUpdate{
					Id:         id,
					Properties: *input,
				}).Execute()
			if err != nil {
				return err
			}

			return printDistribution(c, res)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return addDistributionCreateFlags(cmd)
}

func addDistributionCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagCDNDistributionDomain, "", "", "The public hostname this distribution serves, e.g. cdn.example.com. Must be a valid, unique domain (2-253 chars)")
	cmd.AddStringFlag(constants.FlagCDNDistributionCertificateID, "", "", "Certificate Manager UUID used to terminate HTTPS for the domain. Omit for an HTTP-only distribution")
	cmd.AddStringFlag(constants.FlagCDNDistributionRoutingRules, "", "", "Routing rules as a JSON array (inline string or a path to a .json file). Each rule maps a path prefix + scheme to an upstream origin (host, caching, waf, rateLimitClass, sniMode, geoRestrictions). 1-25 rules. See --routing-rules-example for the format")
	cmd.AddBoolFlag(constants.FlagCDNDistributionRoutingRulesExample, "", false, "Print a ready-to-edit routing-rules JSON template and exit (does not create anything). Redirect to a file, edit it, then pass it via --routing-rules")
	return cmd
}

func setPropertiesFromFlags(c *core.CommandConfig, p *cdn.DistributionProperties) error {
	if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionDomain); viper.IsSet(fn) {
		p.Domain = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionCertificateID); viper.IsSet(fn) {
		p.CertificateId = pointer.From(viper.GetString(fn))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionRoutingRules); viper.IsSet(fn) {
		rr := viper.GetString(fn)
		data, err := getRoutingRulesData(rr)
		if err != nil {
			return fmt.Errorf("error reading routing rules file: %s", err)
		}

		rules, err := getRoutingRulesFromJSON(data)
		if err != nil {
			return fmt.Errorf("error parsing routing rules: %s", err)
		}
		p.RoutingRules = rules
	}

	return nil
}

func getRoutingRulesFromJSON(data []byte) ([]cdn.RoutingRule, error) {
	var rr []cdn.RoutingRule
	err := json.Unmarshal(data, &rr)
	return rr, err
}

func getRoutingRulesData(input string) ([]byte, error) {
	switch _, err := os.Stat(input); {
	case err == nil:
		return os.ReadFile(input)
	case os.IsNotExist(err):
		return []byte(input), nil
	default:
		return nil, err
	}
}

// RoutingRuleExample is the template printed by --routing-rules-example.
// It shows both geoRestrictions variants (allowList vs blockList): supply at
// most one of the two per rule, since the API accepts only one list per rule.
const RoutingRuleExample = `
[
	{
	  "prefix": "/api",
	  "scheme": "http/https",
	  "upstream": {
		"host": "api-origin.example.com",
		"caching": true,
		"waf": true,
		"rateLimitClass": "R500",
		"sniMode": "distribution",
		"geoRestrictions": {
		  "allowList": ["DE", "US"]
		}
	  }
	},
	{
	  "prefix": "/static",
	  "scheme": "http/https",
	  "upstream": {
		"host": "static-origin.example.com",
		"caching": false,
		"waf": false,
		"rateLimitClass": "R10",
		"sniMode": "origin",
		"geoRestrictions": {
		  "blockList": ["CN", "RU"]
		}
	  }
	}
]
`
