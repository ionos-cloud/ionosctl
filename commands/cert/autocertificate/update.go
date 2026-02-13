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
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
)

func AutocertificatePutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update an AutoCertificate.",
		Example:   "ionosctl certmanager autocertificate update --autocertificate-id ID --name NEWNAME",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagAutocertificateID, constants.FlagName); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			autocertificateId, _ := c.Command.Command.Flags().GetString(constants.FlagAutocertificateID)
			g, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesFindById(context.Background(), autocertificateId).Execute()
			if err != nil {
				return err
			}
			return UpdateAutocertificatePrint(c, g)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagAutocertificateID, constants.FlagIdShort, "", "Provide the specified AutoCertificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return AutocertificateIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new name of the AutoCertificate", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func UpdateAutocertificatePrint(c *core.CommandConfig, r cert.AutoCertificateRead) error {
	input := r.Properties

	if c.Command.Command.Flags().Changed(constants.FlagName) {
		input.Name, _ = c.Command.Command.Flags().GetString(constants.FlagName)
	}

	autocertificateid, _ := c.Command.Command.Flags().GetString(constants.FlagAutocertificateID)
	rn, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesPatch(context.Background(), autocertificateid).
		AutoCertificatePatch(cert.AutoCertificatePatch{
			Properties: cert.PatchName{Name: input.Name},
		}).Execute()

	if err != nil {
		return fmt.Errorf("failed to update the AutoCertificate's name: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.CertManagerAutocertificate, rn,
		tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return fmt.Errorf("failed generating the output: %w", err)
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
