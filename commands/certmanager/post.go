package certmanager

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/cobra"
)

var RequiredFlagSets = [16][]string{{"certificate-name", "certificate", "certificate-chain", "private-key"},
	{"certificate-name-path", "certificate", "certificate-chain", "private-key"},
	{"certificate-name", "certificate-path", "certificate-chain", "private-key"},
	{"certificate-name-path", "certificate-path", "certificate-chain", "private-key"},
	{"certificate-name", "certificate", "certificate-chain-path", "private-key"},
	{"certificate-name-path", "certificate", "certificate-chain-path", "private-key"},
	{"certificate-name", "certificate-path", "certificate-chain-path", "private-key"},
	{"certificate-name-path", "certificate-path", "certificate-chain-path", "private-key"},
	{"certificate-name", "certificate", "certificate-chain", "private-key-path"},
	{"certificate-name-path", "certificate", "certificate-chain", "private-key-path"},
	{"certificate-name", "certificate-path", "certificate-chain", "private-key-path"},
	{"certificate-name-path", "certificate-path", "certificate-chain", "private-key-path"},
	{"certificate-name", "certificate", "certificate-chain-path", "private-key-path"},
	{"certificate-name-path", "certificate", "certificate-chain-path", "private-key-path"},
	{"certificate-name", "certificate-path", "certificate-chain-path", "private-key-path"},
	{"certificate-name-path", "certificate-path", "certificate-chain-path", "private-key-path"}}

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
			var err error
			for _, RequiredFlagSet := range RequiredFlagSets {
				err = core.CheckRequiredFlagsSets(c.Command, c.NS, RequiredFlagSet)
				if err == nil {
					return nil
				}
			}
			return err
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Adding Certificate...")
			var name, certificate, certificateChain, privateKey string

			name, err := SetProperties(c, "certificate-name")
			if err != nil {
				return err
			}
			certificate, err = SetProperties(c, "certificate")
			if err != nil {
				return err
			}
			certificateChain, err = SetProperties(c, "certificate-chain")
			if err != nil {
				return err
			}
			privateKey, err = SetProperties(c, "private-key")
			if err != nil {
				return err
			}

			createProperties.SetName(name)
			createProperties.SetCertificate(certificate)
			createProperties.SetCertificateChain(certificateChain)
			createProperties.SetPrivateKey(privateKey)

			Dto := sdkgo.NewCertificatePostDtoWithDefaults()

			Dto.SetProperties(createProperties)

			cert, _, err := c.CertificateManagerServices.Certs().Post(*Dto)
			if err != nil {
				return err
			}

			return c.Printer.Print(getCertPrint(nil, c, &[]sdkgo.CertificateDto{cert}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag("certificate-name", "", "", "Specify name of the certificate (required)")
	cmd.AddStringFlag("certificate", "", "", "Specify the certificate itself (required)")
	cmd.AddStringFlag("certificate-chain", "", "", "Specify the certificate chain (required)")
	cmd.AddStringFlag("private-key", "", "", "Specify the private key (required)")

	cmd.AddStringFlag("certificate-name-path", "", "", "Specify name of the certificate from a file (required)")
	cmd.AddStringFlag("certificate-path", "", "", "Specify the certificate itself from a file (required)")
	cmd.AddStringFlag("certificate-chain-path", "", "", "Specify the certificate chain from a file (required)")
	cmd.AddStringFlag("private-key-path", "", "", "Specify the private key from a file (required)")

	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func SetProperties(c *core.CommandConfig, property string) (string, error) {
	propertyValue, err := c.Command.Command.Flags().GetString(property)
	if propertyValue == "" {
		if err != nil {
			return "", err
		}
		propertyPathName := fmt.Sprintf("%s-path", property)
		propertyPath, err := c.Command.Command.Flags().GetString(propertyPathName)
		if err != nil {
			return "", err
		}
		propertyBytes, err := os.ReadFile(propertyPath)
		if err != nil {
			return "", err
		}
		propertyValue = string(propertyBytes)
		return propertyValue, nil
	}
	return propertyValue, nil
}

// func ifEmptyMarkRequired(c *core.PreCommandConfig, flagName string) error {
// 	flagNamePath := fmt.Sprintf("%s-path", flagName)
// 	if c.Command.Command.Flag(flagName).Value.String() == "" && c.Command.Command.Flag(flagNamePath).Value.String() == "" {
// 		err := c.Command.Command.MarkFlagRequired(flagName)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
