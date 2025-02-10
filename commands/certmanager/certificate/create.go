package certificate

import (
	"context"
	"fmt"
	"os"

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
			err := c.Command.Command.MarkFlagRequired(certmanager.FlagCertName)
			if err != nil {
				panic("unable to mark flag " + certmanager.FlagCertName + " as required. " +
					"Please open an issue at https://github.com/ionos-cloud/ionosctl/issues/new")
			}

			c.Command.Command.MarkFlagsMutuallyExclusive(certmanager.FlagCert, certmanager.FlagCertPath)
			c.Command.Command.MarkFlagsMutuallyExclusive(certmanager.FlagCertChain, certmanager.FlagCertChainPath)
			c.Command.Command.MarkFlagsMutuallyExclusive(certmanager.FlagPrivateKey, certmanager.FlagPrivateKeyPath)

			c.Command.Command.MarkFlagsOneRequired(certmanager.FlagCert, certmanager.FlagCertPath)
			c.Command.Command.MarkFlagsOneRequired(certmanager.FlagCertChain, certmanager.FlagCertChainPath)
			c.Command.Command.MarkFlagsOneRequired(certmanager.FlagPrivateKey, certmanager.FlagPrivateKeyPath)

			return nil
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
	if fn := core.GetFlagName(c.NS, certmanager.FlagCert); viper.IsSet(fn) {
		input.Certificate = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, certmanager.FlagCertChain); viper.IsSet(fn) {
		input.CertificateChain = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, certmanager.FlagPrivateKey); viper.IsSet(fn) {
		input.PrivateKey = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, certmanager.FlagCertPath); viper.IsSet(fn) {
		bytes, err := os.ReadFile(viper.GetString(fn))
		if err != nil {
			return fmt.Errorf("error reading certificate file: %w", err)
		}
		input.Certificate = string(bytes)
	}
	if fn := core.GetFlagName(c.NS, certmanager.FlagCertChainPath); viper.IsSet(fn) {
		bytes, err := os.ReadFile(viper.GetString(fn))
		if err != nil {
			return fmt.Errorf("error reading certificate chain file: %w", err)
		}
		input.CertificateChain = string(bytes)
	}
	if fn := core.GetFlagName(c.NS, certmanager.FlagPrivateKeyPath); viper.IsSet(fn) {
		bytes, err := os.ReadFile(viper.GetString(fn))
		if err != nil {
			return fmt.Errorf("error reading private key file: %w", err)
		}
		input.PrivateKey = string(bytes)
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
