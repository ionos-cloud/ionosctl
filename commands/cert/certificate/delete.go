package certificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CertDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "certmanager",
			Resource:  "certificates",
			Verb:      "delete",
			Aliases:   []string{"d"},
			ShortDesc: "Delete Certificate by ID or all Certificates",
			LongDesc:  "Use this command to delete a Certificate by ID.",
			Example: "ionosctl certificate-manager delete " +
				core.FlagsUsage(constants.FlagCertId) + "\n" +
				"ionosctl certificate-manager delete --all",
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagCertId, constants.FlagIdShort, "", "Response delete a single certificate (required)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCertId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return CertificatesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Response delete all certificates")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(defaultCertificateCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return defaultCertificateCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	var err error

	allFlag, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
	if err != nil {
		return err
	}

	if allFlag {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all Certificates..."))

		certs, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesGet(context.Background()).Execute()
		if err != nil {
			return err
		}

		for _, cert := range certs.Items {
			msg := fmt.Sprintf("delete Certificate ID: %s", cert.Id)
			if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().CertManagerClient.CertificateApi.CertificatesDelete(context.Background(), cert.Id).Execute()
			if err != nil {
				return err
			}
		}
	} else {
		id, err := c.Command.Command.Flags().GetString(constants.FlagCertId)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("delete Certificate ID: %s", id)
		if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = client.Must().CertManagerClient.CertificateApi.CertificatesDelete(context.Background(), id).Execute()

		return err
	}
	return err
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagCertId},
		[]string{constants.ArgAll},
	)
}
