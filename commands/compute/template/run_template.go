package template

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTemplateId)
}

func RunTemplateList(c *core.CommandConfig) error {

	templates, resp, err := c.CloudApiV6Services.Templates().List()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	templatesConverted, err := resource2table.ConvertTemplatesToTable(templates.Templates)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(templates, templatesConverted,
		tabheaders.GetHeadersAllDefault(defaultTemplateCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunTemplateGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Template with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId))))

	tpl, resp, err := c.CloudApiV6Services.Templates().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	templateConverted, err := resource2table.ConvertTemplateToTable(tpl.Template)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(tpl, templateConverted,
		tabheaders.GetHeadersAllDefault(defaultTemplateCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
