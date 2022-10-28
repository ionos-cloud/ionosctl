package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/cobra"
)

var createProperties = sdkgo.CertificatePostPropertiesDto{}

func CertPostCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificates",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Add a new Certificate",
		LongDesc:  "Use this command to add a Certificate.",
		Example:   "ionosctl certificate-manager create",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired("certificate-name")
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired("certificate")
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired("certificate-chain")
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired("private-key")
			if err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			Dto := sdkgo.NewCertificatePostDtoWithDefaults()
			name, err := c.Command.Command.Flags().GetString("certificate-name")
			if err != nil {
				return err
			}
			certificate, err := c.Command.Command.Flags().GetString("certificate")
			if err != nil {
				return err
			}
			certificateChain, err := c.Command.Command.Flags().GetString("certificate-chain")
			if err != nil {
				return err
			}
			privateKey, err := c.Command.Command.Flags().GetString("private-key")
			if err != nil {
				return err
			}

			createProperties.SetName(name)
			createProperties.SetCertificate(certificate)
			createProperties.SetCertificateChain(certificateChain)
			createProperties.SetPrivateKey(privateKey)

			Dto.SetProperties(createProperties)

			c.Printer.Verbose("Adding Certificate...")
			cert, r, err := c.CertificateManagerServices.Certs().Post(*Dto)
			if err != nil {
				return err
			}
			var list = sdkgo.CertificateCollectionDto{}

			*list.GetItems() = append(*list.GetItems(), cert)
			return c.Printer.Print(getCertPrint(r, c, &list))
		},
		InitClient: true,
	})

	cmd.AddStringFlag("certificate-name", "", "", "Response get a single certificate (required)")
	cmd.AddStringFlag("certificate", "", "", "Response get a single certificate (required)")
	cmd.AddStringFlag("certificate-chain", "", "", "Response get a single certificate (required)")
	cmd.AddStringFlag("private-key", "", "", "Response get a single certificate (required)")
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
