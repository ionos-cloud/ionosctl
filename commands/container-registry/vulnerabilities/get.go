package vulnerabilities

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
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

	cmd.Command.Flags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(constants.FlagVulnerabilityId, "", "", "Vulnerability ID")

	return cmd
}

func PreCmdGet(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagVulnerabilityId)
}

func CmdGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)
	vulnId := viper.GetString(core.GetFlagName(c.NS, constants.FlagVulnerabilityId))

	vulnerability, _, err := client.Must().RegistryClient.VulnerabilitiesApi.VulnerabilitiesFindByID(
		context.
			Background(), vulnId,
	).Execute()
	if err != nil {
		return err
	}

	vulnerabilityConverted, err := resource2table.ConvertContainerRegistryVulnerabilityToTable(vulnerability)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		vulnerability, vulnerabilityConverted, tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
