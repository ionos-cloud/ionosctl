package lan

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunDcNatGatewayLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgLanId)
}

func PreRunDcNatGatewayLanRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgLanId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgAll},
	)
}

func RunNatGatewayLanList(c *core.CommandConfig) error {
	ng, resp, err := c.CloudApiV6Services.NatGateways().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("properties.lans", jsonpaths.NatGatewayLan, ng.NatGateway,
		tabheaders.GetHeadersAllDefault(defaultNatGatewayLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunNatGatewayLanAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))

	ng, _, err := c.CloudApiV6Services.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Adding NatGateway with id %v to Datacenter with id: %v", natGatewayId, dcId))

	input := getNewNatGatewayLanInfo(c, ng)
	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(dcId, natGatewayId, *input)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("properties.lans", jsonpaths.NatGatewayLan, ng.NatGateway,
		tabheaders.GetHeadersAllDefault(defaultNatGatewayLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunNatGatewayLanRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllNatGatewayLans(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove nat gateway lan", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))

	ng, _, err := c.CloudApiV6Services.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Removing NatGateway with id %v to Datacenter with id: %v", natGatewayId, dcId))

	input := removeNatGatewayLanInfo(c, ng)
	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(dcId, natGatewayId, *input)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("NAT Gateway Lan successfully deleted"))
	return nil
}

func RemoveAllNatGatewayLans(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("NatGateway ID: %v", natGatewayId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting NatGateway..."))

	natGateway, resp, err := c.CloudApiV6Services.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}

	natGatewayProperties, ok := natGateway.GetPropertiesOk()
	if !ok || natGatewayProperties == nil {
		return fmt.Errorf("could not get NAT Gateway properties")
	}

	lansOk, ok := natGatewayProperties.GetLansOk()
	if !ok || lansOk == nil {
		return fmt.Errorf("could not get items of NAT Gateway Lans")
	}

	if len(*lansOk) <= 0 {
		return fmt.Errorf("no NAT Gateway Lans found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("NAT Gateway Lans to be removed:"))
	for _, lan := range *lansOk {
		if id, ok := lan.GetIdOk(); ok && id != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("NAT Gateway Lan Id: %v", string(*id)))
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove all the NAT Gateways Lans", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing all the NAT Gateways Lans..."))

	proper := make([]ionoscloud.NatGatewayLanProperties, 0)
	if natGateway != nil {
		if properties, ok := natGateway.GetPropertiesOk(); ok && properties != nil {
			natGatewaysProps := &resources.NatGatewayProperties{
				NatGatewayProperties: ionoscloud.NatGatewayProperties{
					Lans: &proper,
				},
			}

			natGateway, resp, err = c.CloudApiV6Services.NatGateways().Update(dcId, natGatewayId, *natGatewaysProps)
			if resp != nil && request.GetId(resp) != "" {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
			}
			if err != nil {
				return err
			}

			if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
				return err
			}
		}
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("NAT Gateway Lans successfully deleted"))
	return nil
}

func getNewNatGatewayLanInfo(c *core.CommandConfig, oldNg *resources.NatGateway) *resources.NatGatewayProperties {
	var proper []ionoscloud.NatGatewayLanProperties

	if oldNg != nil {
		if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				proper = *lans
			}
		}
	}

	input := ionoscloud.NatGatewayLanProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
		lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
		input.SetId(lanId)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Id set: %v", lanId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		gatewayIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps))
		input.SetGatewayIps(gatewayIps)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property GatewayIps set: %v", gatewayIps))
	}

	proper = append(proper, input)

	return &resources.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &proper,
		},
	}
}

func removeNatGatewayLanInfo(c *core.CommandConfig, oldNg *resources.NatGateway) *resources.NatGatewayProperties {
	proper := make([]ionoscloud.NatGatewayLanProperties, 0)

	if oldNg != nil {
		if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					if id, ok := lanItem.GetIdOk(); ok && id != nil {
						if *id != viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
							proper = append(proper, lanItem)
						}
					}
				}
			}
		}
	}

	return &resources.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &proper,
		},
	}
}

func PreRunDcNatGatewayIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId)
}
