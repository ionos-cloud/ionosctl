package vulnerabilities

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func VulnerabilitiesGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "vulnerabilities",
			Verb:       "get",
			ShortDesc:  "Retrieve a vulnerability",
			LongDesc:   "Retrieve a vulnerability",
			Example:    "ionosctl container-registry vulnerabilities get",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagVulnerabilityId, "", "", "Vulnerability ID")

	return cmd
}

func PreCmdGet(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagVulnerabilityId)
}

func CmdGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	vulnId := viper.GetString(core.GetFlagName(c.NS, constants.FlagVulnerabilityId))

	vulnerability, _, err := client.Must().RegistryClient.VulnerabilitiesApi.VulnerabilitiesFindByID(
		context.
			Background(), vulnId,
	).Execute()
	if err != nil {
		return err
	}

	return c.Out(table.Sprint(allCols, vulnerability, cols))
}
