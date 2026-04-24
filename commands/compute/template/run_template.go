package template

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTemplateId)
}

func RunTemplateList(c *core.CommandConfig) error {
	templates, resp, err := c.CloudApiV6Services.Templates().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allTemplateCols).Prefix("items").Print(templates.Templates)
}

func RunTemplateGet(c *core.CommandConfig) error {
	c.Verbose("Template with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)))

	tpl, resp, err := c.CloudApiV6Services.Templates().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allTemplateCols).Print(tpl.Template)
}
