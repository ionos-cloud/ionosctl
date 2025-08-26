package autocertificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func AutocertificateFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve an AutoCertificate",
		Example:   "ionosctl certmanager autocertificate get --autocertificate-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagAutocertificateID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			autocertificateId := viper.GetString(core.GetFlagName(c.NS, constants.FlagAutocertificateID))
			r, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesFindById(context.Background(), autocertificateId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the AutoCertificate: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.CertManagerAutocertificate, r, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the output: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagAutocertificateID, constants.FlagIdShort, "", "Provide the specified AutoCertificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return AutocertificateIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
