package provider

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

func ProviderPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "provider",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Modify an Provider",
		Example:   "ionosctl certmanager provider update --provider-id ID --name NEWNAME",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagProviderID, constants.FlagName); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			providerId, _ := c.Command.Command.Flags().GetString(constants.FlagProviderID)
			g, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersFindById(context.Background(), providerId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the Provider: %w", err)
			}
			return UpdateProviderPrint(c, g)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagProviderID, constants.FlagIdShort, "", "The certificate Provider used to issue the certificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return ProviderIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new name of the Provider", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func UpdateProviderPrint(c *core.CommandConfig, r cert.ProviderRead) error {
	input := r.Properties

	if c.Command.Command.Flags().Changed(constants.FlagName) {
		input.Name, _ = c.Command.Command.Flags().GetString(constants.FlagName)
	}

	providerid, _ := c.Command.Command.Flags().GetString(constants.FlagProviderID)
	rn, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersPatch(context.Background(), providerid).
		ProviderPatch(cert.ProviderPatch{
			Properties: cert.PatchName{Name: input.Name},
		}).Execute()

	if err != nil {
		return fmt.Errorf("failed to update the Provider's name: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.CertManagerProvider, rn,
		tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return fmt.Errorf("failed generating the output: %w", err)
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
