package certmanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/viper"
	"os"
)

var createProperties = sdkgo.CertificatePostPropertiesDto{}

func CertCreateCmd() *core.Command {
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

	cmd.AddStringFlag(FlagCertName, "n", "", "Specify name of the certificate", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagCert, "", "", "Specify the certificate itself (required either this or --certificate-path)")
	cmd.AddStringFlag(FlagCertChain, "", "", "Specify the certificate chain (required either this or --certificate-chain-path)")
	cmd.AddStringFlag(FlagPrivateKey, "", "", "Specify the private key (required either this or --private-key-path)")

	cmd.AddStringFlag(FlagCertPath, "", "", "Specify the certificate itself from a file (required either this or --certificate)")
	cmd.AddStringFlag(FlagCertChainPath, "", "", "Specify the certificate chain from a file (required either this or --certificate-chain)")
	cmd.AddStringFlag(FlagPrivateKeyPath, "", "", "Specify the private key from a file (required either this or --private-key)")

	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "Get response without headers")

	return cmd
}

func GetPropertyWithFallback(c *core.CommandConfig, property string, propertyPath string) (string, error) {
	propertyValue, err := c.Command.Command.Flags().GetString(property)
	if err != nil {
		return "", err
	}
	if propertyValue != "" {
		return propertyValue, nil
	}
	propertyValuePath, err := c.Command.Command.Flags().GetString(propertyPath)
	if err != nil {
		return "", err
	}
	propertyBytes, err := os.ReadFile(propertyValuePath)
	if err != nil {
		return "", err
	}
	return string(propertyBytes), nil
}

func CmdPost(c *core.CommandConfig) error {
	c.Printer.Verbose("Adding Certificate...")
	var name, certificate, certificateChain, privateKey string

	name, err := c.Command.Command.Flags().GetString(FlagCertName)
	if err != nil {
		return err
	}
	certificate, err = GetPropertyWithFallback(c, FlagCert, FlagCertPath)
	if err != nil {
		return err
	}
	certificateChain, err = GetPropertyWithFallback(c, FlagCertChain, FlagCertChainPath)
	if err != nil {
		return err
	}
	privateKey, err = GetPropertyWithFallback(c, FlagPrivateKey, FlagPrivateKeyPath)
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

func PreCmdPost(c *core.PreCommandConfig) error {
	err := errors.New("")
	if !viper.IsSet(core.GetFlagName(c.NS, FlagCert)) && !viper.IsSet(core.GetFlagName(c.NS, FlagCertPath)) {
		err = fmt.Errorf("%veither --%s or --%s must be set\n", err, FlagCert, FlagCertPath)
	}
	if !viper.IsSet(core.GetFlagName(c.NS, FlagCertChain)) && !viper.IsSet(core.GetFlagName(c.NS, FlagCertChainPath)) {
		err = fmt.Errorf("%veither --%s or --%s must be set\n", err, FlagCertChain, FlagCertChainPath)
	}
	if !viper.IsSet(core.GetFlagName(c.NS, FlagPrivateKey)) && !viper.IsSet(core.GetFlagName(c.NS, FlagPrivateKeyPath)) {
		err = fmt.Errorf("%veither --%s or --%s must be set\n", err, FlagPrivateKey, FlagPrivateKeyPath)
	}
	err = c.Command.Command.MarkFlagRequired(FlagCertName)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("%v\nUsage:\n%s\n\nFor more details, see '%s --help'.",
			err,
			UsageCert,
			c.Command.CommandPath(),
		)
	}
	return err
}
