package contract

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func RunContractGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Contract with resource limits: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits))))

	contractResource, resp, err := c.CloudApiV6Services.Contracts().Get()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	var out string

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits)) {
		switch strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceLimits))) {
		case "CORES":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractCoresCols)
		case "RAM":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractRamCols)
		case "HDD":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractHddCols)
		case "SSD":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractSsdCols)
		case "DAS":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractDasCols)
		case "IPS":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractIpsCols)
		case "K8S":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractK8sCols)
		case "NLB":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractNlbCols)
		case "NAT":
			out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts, contractNatCols)
		}
	} else {
		out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Contract, contractResource.Contracts,
			tabheaders.GetHeaders(allContractCols, defaultContractCols, cols))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
