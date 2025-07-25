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
		ShortDesc: "Create a CDN distribution. Wiki: https://docs.ionos.com/cloud/network-services/cdn/dcd-how-tos/create-cdn-distribution",
		Example:   "ionosctl cdn ds create --domain foo-bar.com --certificate-id id --routing-rules rules.json",
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
	cmd.AddStringFlag(constants.FlagCDNDistributionDomain, "", "", "The domain of the distribution")
	cmd.AddStringFlag(constants.FlagCDNDistributionCertificateID, "", "", "The ID of the certificate")
	cmd.AddStringFlag(constants.FlagCDNDistributionRoutingRules, "", "", "The routing rules of the distribution. JSON string or file path of routing rules")
	cmd.AddBoolFlag(constants.FlagCDNDistributionRoutingRulesExample, "", false, "Print an example of routing rules")
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

const RoutingRuleExample = `
[
	{
	  "prefix": "/api",
	  "scheme": "http/https",
	  "upstream": {
		"caching": true,
		"geoRestrictions": {
		  "allowList": ["CN", "RU"]
		},
		"host": "clitest.example.com",
		"rateLimitClass": "R500",
		"sniMode": "distribution",
		"waf": true
	  }
	},
    {
	  "prefix": "/api2",
	  "scheme": "http/https",
	  "upstream": {
		"caching": false,
		"geoRestrictions": {
		  "blockList": ["CN", "RU"]
		},
		"host": "server2.example.com",
		"rateLimitClass": "R10",
		"sniMode": "origin",
		"waf": false
      }
   }
]
`
