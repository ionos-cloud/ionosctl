package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func CertDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificates",
		Verb:      "delete",
		Aliases:   []string{"g"},
		ShortDesc: "Delete Certificate by ID",
		LongDesc:  "Use this command to delete a Certificate by ID.",
		Example:   "ionsoclt certificate-manager delete --certificate-id 12345678-1234-1234-1234-123456789012",
		PreCmdRun: PreCmdDelete,
		CmdRun: CmdDelete,
		InitClient: true,
	})

	cmd.AddStringFlag("certificate-id", "", "", "Response get a single certificate (required)")
	cmd.AddBoolFlag("all", "", false, "Response get a single certificate")
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	var err error
	allFlag := c.Command.Command.Flag("all").Changed
	if allFlag {
		c.Printer.Verbose("Deleting all Certificates...")
		certs, _, err := c.CertificateManagerServices.Certs().List()
		if err != nil {
			return err
		}

		for _, cert := range *certs.Items {
			_, err = c.CertificateManagerServices.Certs().Delete(*cert.Id)
			if err != nil {
				return err
			}
		}

	} else {
		c.Printer.Verbose("Deleting Certificate...")
		id, err := c.Command.Command.Flags().GetString("certificate-id")
		if err != nil {
			return err
		}
		_, err = c.CertificateManagerServices.Certs().Delete(id)

		return err
	}
	return err
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{"certificate-id"},
		[]string{"all"},
	)
}
