package autocertificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func AutocertificateFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get an auto-certificate by ID",
		LongDesc:  "Retrieve a single auto-certificate by ID, showing its provider, common name, additional names, key algorithm, and provisioning state (PROVISIONING, AVAILABLE, or FAILED).",
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

			return c.Printer(allCols).Print(r)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagAutocertificateID, constants.FlagIdShort, "", "The ID (UUID) of the auto-certificate to retrieve", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return AutocertificateIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
