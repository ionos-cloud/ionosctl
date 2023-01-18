package certmanager

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/cobra"
)

func CertDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete Certificate by ID or all Certificates",
		LongDesc:   "Use this command to delete a Certificate by ID.",
		Example:    "ionsoctl certificate-manager delete --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422",
		PreCmdRun:  PreCmdDelete,
		CmdRun:     CmdDelete,
		InitClient: true,
	})

	cmd.AddStringFlag(FlagCertId, constants.FlagIdP, "", "Response delete a single certificate (required)")
	_ = cmd.Command.RegisterFlagCompletionFunc(FlagCertId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return CertificatesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Response delete all certificates")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	var err error
	allFlag, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
	if err != nil {
		return err
	}
	if allFlag {
		c.Printer.Verbose("Deleting all Certificates...")
		certs, _, err := c.CertificateManagerServices.Certs().List()
		if err != nil {
			return err
		}
		for _, cert := range *certs.Items {
			msg := fmt.Sprintf("delete Certificate ID: %s", *cert.Id)
			if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
				return err
			}
			_, err = c.CertificateManagerServices.Certs().Delete(*cert.Id)
			if err != nil {
				return err
			}
		}

	} else {
		id, err := c.Command.Command.Flags().GetString(FlagCertId)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("delete Certificate ID: %s", id)
		if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
			return err
		}
		_, err = c.CertificateManagerServices.Certs().Delete(id)

		return err
	}
	return err
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{FlagCertId},
		[]string{constants.ArgAll},
	)
}
