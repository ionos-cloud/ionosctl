package certificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
)

func CertUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "update",
		Aliases:    []string{"u"},
		ShortDesc:  "Update Certificate name",
		LongDesc:   "Use this change a certificate's name.",
		Example:    "ionosctl certificate-manager update --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422",
		PreCmdRun:  PreCmdPatch,
		CmdRun:     CmdPatch,
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCertId, constants.FlagIdShort, "", "Provide the specified Certificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return CertificatesIds()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)
	cmd.AddStringFlag(constants.FlagCertName, "n", "", "Provide new certificate name", core.RequiredFlagOption())

	return cmd
}

func CmdPatch(c *core.CommandConfig) error {
	id, err := c.Command.Command.Flags().GetString(constants.FlagCertId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(constants.FlagCertName)
	if err != nil {
		return err
	}

	cert, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesPatch(context.Background(), id).
		CertificatePatch(
			cert.CertificatePatch{
				Properties: cert.PatchName{
					Name: name,
				},
			},
		).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(cert)
}

func PreCmdPatch(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(constants.FlagCertId)
	if err != nil {
		return err
	}

	err = c.Command.Command.MarkFlagRequired(constants.FlagCertName)
	if err != nil {
		return err
	}

	return nil
}
