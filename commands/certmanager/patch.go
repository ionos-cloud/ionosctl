package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"

	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
)

var editProperties = sdkgo.CertificatePatchPropertiesDto{}

func CertPatchCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "patch",
		Aliases:    []string{"p"},
		ShortDesc:  "Get Certificate by ID",
		LongDesc:   "Use this command to retrieve a Certificate by ID.",
		Example:    "ionosctl certificate-manager get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422",
		PreCmdRun:  PreCmdPatch,
		CmdRun:     CmdPatch,
		InitClient: true,
	})

	cmd.AddStringFlag(CertId, "i", "", "Provide certificate ID (required)")
	cmd.AddStringFlag(CertName, "n", "", "Provide new certificate name (required)")
	_ = cmd.Command.RegisterFlagCompletionFunc(ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdPatch(c *core.CommandConfig) error {
	c.Printer.Verbose("Patching Certificate...")
	id, err := c.Command.Command.Flags().GetString(CertId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(CertName)
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
	err := c.Command.Command.MarkFlagRequired(CertId)
	if err != nil {
		return err
	}
	err = c.Command.Command.MarkFlagRequired(CertName)
	if err != nil {
		return err
	}
	return nil
}
