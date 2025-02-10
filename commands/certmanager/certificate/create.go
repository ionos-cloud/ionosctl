package certificate

import (
	"context"
	"fmt"
	"os"

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
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificate",
		Verb:      "create",
		Aliases:   []string{"add", "a", "c", "post"},
		ShortDesc: "Add a new Certificate",
		LongDesc:  "Use this command to add a Certificate.",
		Example: "ionosctl certmanager certificate create " + core.FlagsUsage(constants.FlagCertName,
			constants.FlagCertChainPath, constants.FlagCertPath, constants.FlagPrivateKeyPath),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagCertName)
			if err != nil {
				panic("unable to mark flag " + constants.FlagCertName + " as required. " +
					"Please open an issue at https://github.com/ionos-cloud/ionosctl/issues/new")
			}

			c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagCert, constants.FlagCertPath)
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagCertChain, constants.FlagCertChainPath)
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagPrivateKey, constants.FlagPrivateKeyPath)

			c.Command.Command.MarkFlagsOneRequired(constants.FlagCert, constants.FlagCertPath)
			c.Command.Command.MarkFlagsOneRequired(constants.FlagCertChain, constants.FlagCertChainPath)
			c.Command.Command.MarkFlagsOneRequired(constants.FlagPrivateKey, constants.FlagPrivateKeyPath)

			return nil
		},
		CmdRun:     CmdPost,
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCertName, "n", "", "Specify name of the certificate", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagCert, "", "", "Specify the certificate itself (required either this or --certificate-path)")
	cmd.AddStringFlag(constants.FlagCertChain, "", "", "Specify the certificate chain (required either this or --certificate-chain-path)")
	cmd.AddStringFlag(constants.FlagPrivateKey, "", "", "Specify the private key (required either this or --private-key-path)")

	cmd.AddStringFlag(constants.FlagCertPath, "", "", "Specify the certificate itself from a file (required either this or --certificate)")
	cmd.AddStringFlag(constants.FlagCertChainPath, "", "", "Specify the certificate chain from a file (required either this or --certificate-chain)")
	cmd.AddStringFlag(constants.FlagPrivateKeyPath, "", "", "Specify the private key from a file (required either this or --private-key)")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(defaultCertificateCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultCertificateCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdPost(c *core.CommandConfig) error {
	input := cert.Certificate{}

	if fn := core.GetFlagName(c.NS, constants.FlagCertName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagCert); viper.IsSet(fn) {
		input.Certificate = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagCertChain); viper.IsSet(fn) {
		input.CertificateChain = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPrivateKey); viper.IsSet(fn) {
		input.PrivateKey = viper.GetString(fn)
	}

	if fn := core.GetFlagName(c.NS, constants.FlagCertPath); viper.IsSet(fn) {
		bytes, err := os.ReadFile(viper.GetString(fn))
		if err != nil {
			return fmt.Errorf("error reading certificate file: %w", err)
		}
		input.Certificate = string(bytes)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagCertChainPath); viper.IsSet(fn) {
		bytes, err := os.ReadFile(viper.GetString(fn))
		if err != nil {
			return fmt.Errorf("error reading certificate chain file: %w", err)
		}
		input.CertificateChain = string(bytes)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPrivateKeyPath); viper.IsSet(fn) {
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
