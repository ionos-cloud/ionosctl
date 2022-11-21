package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/core"

	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
)

var editProperties = sdkgo.CertificatePatchPropertiesDto{}

func CertUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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

	cmd.AddStringFlag(FlagCertId, "i", "", "Provide certificate ID", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagCertName, "n", "", "Provide new certificate name", core.RequiredFlagOption())

	return cmd
}

func CmdPatch(c *core.CommandConfig) error {
	c.Printer.Verbose("Patching Certificate...")
	id, err := c.Command.Command.Flags().GetString(FlagCertId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(FlagCertName)
	if err != nil {
		return err
	}

	editProperties.SetName(name)

	input := sdkgo.NewCertificatePatchDto(editProperties)
	cert, _, err := c.CertificateManagerServices.Certs().Patch(id, *input)
	if err != nil {
		return err
	}

	return c.Printer.Print(getCertPrint(nil, c, &[]sdkgo.CertificateDto{cert}))
}

func PreCmdPatch(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(FlagCertId)
	if err != nil {
		return err
	}
	err = c.Command.Command.MarkFlagRequired(FlagCertName)
	if err != nil {
		return err
	}
	return nil
}
