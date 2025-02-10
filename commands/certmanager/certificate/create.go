package certificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/certmanager"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CertCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificate",
		Verb:      "create",
		Aliases:   []string{"add", "a", "c", "post"},
		ShortDesc: "Add a new Certificate",
		LongDesc:  "Use this command to add a Certificate.",
		Example:   "TODO",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			var err error
			if !viper.IsSet(core.GetFlagName(c.NS, certmanager.FlagCert)) && !viper.IsSet(core.GetFlagName(c.NS, certmanager.FlagCertPath)) {
				err = fmt.Errorf("%veither --%s or --%s must be set", err, certmanager.FlagCert, certmanager.FlagCertPath)
			}
			if !viper.IsSet(core.GetFlagName(c.NS, certmanager.FlagCertChain)) && !viper.IsSet(core.GetFlagName(c.NS, certmanager.FlagCertChainPath)) {
				err = fmt.Errorf("%veither --%s or --%s must be set", err, certmanager.FlagCertChain, certmanager.FlagCertChainPath)
			}
			if !viper.IsSet(core.GetFlagName(c.NS, certmanager.FlagPrivateKey)) && !viper.IsSet(core.GetFlagName(c.NS, certmanager.FlagPrivateKeyPath)) {
				err = fmt.Errorf("%veither --%s or --%s must be set", err, certmanager.FlagPrivateKey, certmanager.FlagPrivateKeyPath)
			}
			err = c.Command.Command.MarkFlagRequired(certmanager.FlagCertName)

			return err
		},
		CmdRun:     CmdPost,
		InitClient: true,
	})

	cmd.AddStringFlag(certmanager.FlagCertName, "n", "", "Specify name of the certificate", core.RequiredFlagOption())
	cmd.AddStringFlag(certmanager.FlagCert, "", "", "Specify the certificate itself (required either this or --certificate-path)")
	cmd.AddStringFlag(certmanager.FlagCertChain, "", "", "Specify the certificate chain (required either this or --certificate-chain-path)")
	cmd.AddStringFlag(certmanager.FlagPrivateKey, "", "", "Specify the private key (required either this or --private-key-path)")

	cmd.AddStringFlag(certmanager.FlagCertPath, "", "", "Specify the certificate itself from a file (required either this or --certificate)")
	cmd.AddStringFlag(certmanager.FlagCertChainPath, "", "", "Specify the certificate chain from a file (required either this or --certificate-chain)")
	cmd.AddStringFlag(certmanager.FlagPrivateKeyPath, "", "", "Specify the private key from a file (required either this or --private-key)")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(defaultCertificateCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultCertificateCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdPost(c *core.CommandConfig) error {
	input := cert.Certificate{}

	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}

	cert, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesPost(context.Background()).
		CertificateCreate(cert.CertificateCreate{Properties: input}).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.CertManagerCertificate, cert,
		tabheaders.GetHeadersAllDefault(defaultCertificateCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
