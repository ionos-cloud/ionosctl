package distribution

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"

	cdn "github.com/ionos-cloud/sdk-go-cdn"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagCDNDistributionDomain, constants.FlagCDNDistributionRoutingRules); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := &cdn.DistributionProperties{}
			if err := setPropertiesFromFlags(c, input); err != nil {
				return err
			}

			id := uuidgen.Must()
			res, _, err := client.Must().CDNClient.DistributionsApi.DistributionsPut(context.Background(), id).
				DistributionUpdate(cdn.DistributionUpdate{
					Id:         &id,
					Properties: input,
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
	return cmd
}

func setPropertiesFromFlags(c *core.CommandConfig, p *cdn.DistributionProperties) error {
	if fn := core.GetFlagName(c.NS, constants.FlagCDNDistributionDomain); viper.IsSet(fn) {
		p.Domain = pointer.From(viper.GetString(fn))
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

func getRoutingRulesFromJSON(data []byte) (*[]cdn.RoutingRule, error) {
	var rr []cdn.RoutingRule
	err := json.Unmarshal(data, &rr)
	if err != nil {
		return nil, err
	}
	return &rr, nil
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
