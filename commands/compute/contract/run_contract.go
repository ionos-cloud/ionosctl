package contract

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func RunContractGet(c *core.CommandConfig) error {
	c.Verbose("Contract with resource limits: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)))

	contractResource, resp, err := c.CloudApiV6Services.Contracts().Get()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)) {
		var overrideCols []string
		switch strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits))) {
		case "CORES":
			overrideCols = contractCoresCols
		case "RAM":
			overrideCols = contractRamCols
		case "HDD":
			overrideCols = contractHddCols
		case "SSD":
			overrideCols = contractSsdCols
		case "DAS":
			overrideCols = contractDasCols
		case "IPS":
			overrideCols = contractIpsCols
		case "K8S":
			overrideCols = contractK8sCols
		case "NLB":
			overrideCols = contractNlbCols
		case "NAT":
			overrideCols = contractNatCols
		default:
			return fmt.Errorf("invalid value for --resource-limits: %q. Valid values: CORES, RAM, HDD, SSD, DAS, IPS, K8S, NLB, NAT",
				viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)))
		}
		return c.Out(table.Sprint(allContractCols, contractResource.Contracts, overrideCols, table.WithPrefix("items")))
	}

	return c.Printer(allContractCols).Prefix("items").Print(contractResource.Contracts)
}
