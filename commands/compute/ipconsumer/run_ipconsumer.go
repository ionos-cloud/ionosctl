package ipconsumer

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func RunIpConsumersList(c *core.CommandConfig) error {
	ipBlock, resp, err := c.CloudApiV6Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := ipBlock.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting ip block properties")
	}

	ipCons, ok := properties.GetIpConsumersOk()
	if !ok || ipCons == nil {
		return fmt.Errorf("error getting ip consumers")
	}

	ipsConsumers := make([]ionoscloud.IpConsumer, 0)
	for _, ip := range *ipCons {
		ipsConsumers = append(ipsConsumers, ip)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.IpConsumer, ipsConsumers,
		tabheaders.GetHeaders(allIpConsumerCols, defaultIpConsumerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func PreRunIpBlockId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgIpBlockId)
}
