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

var createProperties = sdkgo.CertificatePostPropertiesDto{}

func CertPostCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "create",
		Aliases:    []string{"c"},
		ShortDesc:  "Add a new Certificate",
		LongDesc:   "Use this command to add a Certificate.",
		Example:    "ionosctl certificate-manager create",
		PreCmdRun:  PreCmdPost,
		CmdRun:     CmdPost,
		InitClient: true,
	})

	cmd.AddStringFlag(CertName, "", "", "Specify name of the certificate (required)")
	cmd.AddStringFlag(Cert, "", "", "Specify the certificate itself (required)")
	cmd.AddStringFlag(CertChain, "", "", "Specify the certificate chain (required)")
	cmd.AddStringFlag(PrivateKey, "", "", "Specify the private key (required)")

	cmd.AddStringFlag(CertNamePath, "", "", "Specify name of the certificate from a file (required)")
	cmd.AddStringFlag(CertPath, "", "", "Specify the certificate itself from a file (required)")
	cmd.AddStringFlag(CertChainPath, "", "", "Specify the certificate chain from a file (required)")
	cmd.AddStringFlag(PrivateKeyPath, "", "", "Specify the private key from a file (required)")

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

func CmdPost(c *core.CommandConfig) error {
	c.Printer.Verbose("Adding Certificate...")
	var name, certificate, certificateChain, privateKey string

	name, err := SetProperties(c, CertName)
	if err != nil {
		return err
	}
	certificate, err = SetProperties(c, Cert)
	if err != nil {
		return err
	}
	certificateChain, err = SetProperties(c, CertChain)
	if err != nil {
		return err
	}
	privateKey, err = SetProperties(c, PrivateKey)
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
}

// couldn't think of a better solution to this
func PreCmdPost(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		RequiredFlagSets[0],
		RequiredFlagSets[1],
		RequiredFlagSets[2],
		RequiredFlagSets[3],
		RequiredFlagSets[4],
		RequiredFlagSets[5],
		RequiredFlagSets[6],
		RequiredFlagSets[7],
		RequiredFlagSets[8],
		RequiredFlagSets[9],
		RequiredFlagSets[10],
		RequiredFlagSets[11],
		RequiredFlagSets[12],
		RequiredFlagSets[13],
		RequiredFlagSets[14],
		RequiredFlagSets[15],
	)

	if err != nil {
		return fmt.Errorf(PostErrorFormatFlag,
			c.Command.CommandPath(),
			PostErrorExample1,
			PostErrorExample2,
			c.Command.CommandPath(),
		)
	}
	return nil
}
