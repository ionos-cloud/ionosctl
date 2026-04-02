package resource

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunResourceType(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagType)
}

func RunResourceList(c *core.CommandConfig) error {
	resourcesListed, resp, err := c.CloudApiV6Services.Users().ListResources()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Resource, resourcesListed.Resources,
		tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunResourceGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Resource with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId))))

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)) {
		resourceListed, resp, err := c.CloudApiV6Services.Users().GetResourceByTypeAndId(
			viper.GetString(core.GetFlagName(c.NS, constants.FlagType)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		)
		if resp != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
		}
		if err != nil {
			return err
		}

		out, err := jsontabwriter.GenerateOutput("", jsonpaths.Resource, resourceListed.Resource,
			tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
		return nil
	}

	resourcesListed, resp, err := c.CloudApiV6Services.Users().GetResourcesByType(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Resource, resourcesListed.Resources,
		tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunGroupResourceList(c *core.CommandConfig) error {

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Listing Resources from Group with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	resourcesListed, resp, err := c.CloudApiV6Services.Groups().ListResources(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Resource, resourcesListed.ResourceGroups,
		tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
