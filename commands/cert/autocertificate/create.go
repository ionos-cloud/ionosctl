package autocertificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cert/provider"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/viper"
)

func AutocertificatePostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Issue an auto-renewing certificate via a provider",
		LongDesc: `Request a certificate that IONOS issues and auto-renews through an ACME provider.

IONOS validates domain ownership with ACME DNS-01 challenges, so the --common-name (and every --subject-alternative-names entry) must belong to a DNS zone hosted in IONOS CLOUD DNS that you manage; IONOS creates the required TXT records automatically. If the matching zone does not exist, issuance fails.

Required: --name, --provider-id (an existing provider, see 'certmanager provider create'), --common-name, and --key-algorithm. Once issued, the certificate renews automatically ~30 days before expiry.`,
		Example: `# Issue a single-domain RSA-2048 certificate
ionosctl certmanager autocertificate create --name web-cert --provider-id PROVIDER_ID --common-name www.example.com --key-algorithm rsa2048

# Issue a certificate covering additional (SAN) domains with a stronger key
ionosctl certmanager autocertificate create --name web-cert --provider-id PROVIDER_ID --common-name www.example.com --key-algorithm rsa4096 --subject-alternative-names app.example.com,api.example.com`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagProviderID, constants.FlagCommonName, constants.FlagKeyAlgorithm)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			input := cert.AutoCertificate{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagProviderID); viper.IsSet(fn) {
				input.Provider = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagCommonName); viper.IsSet(fn) {
				input.CommonName = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagKeyAlgorithm); viper.IsSet(fn) {
				input.KeyAlgorithm = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagSubjectAlternativeNames); viper.IsSet(fn) {
				input.SubjectAlternativeNames = viper.GetStringSlice(fn)
			}

			z, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesPost(context.Background()).
				AutoCertificateCreate(cert.AutoCertificateCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return fmt.Errorf("failed creating the AutoCertificate: %w", err)
			}

			return c.Printer(allCols).Print(z)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "A friendly name to identify the auto-certificate for management purposes. Required")
	cmd.AddStringFlag(constants.FlagProviderID, constants.FlagIdShort, "", "The ID (UUID) of the ACME provider that will issue and renew this certificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return provider.ProviderIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)
	cmd.AddStringFlag(constants.FlagCommonName, "", "", "The primary domain (DNS name) the certificate is issued for, e.g. www.example.com. Must belong to an IONOS CLOUD DNS zone you manage. Required")
	cmd.AddStringFlag(constants.FlagKeyAlgorithm, "", "", "The key algorithm for the generated private key. One of: rsa2048, rsa3072, rsa4096. Required", core.RequiredFlagOption(), core.WithCompletion(
		func() []string {
			return []string{"rsa2048", "rsa3072", "rsa4096"}
		}, constants.CertApiRegionalURL, constants.CertLocations))
	cmd.AddStringSliceFlag(constants.FlagSubjectAlternativeNames, "", []string{""}, "Additional domains (SANs) to cover with the same certificate, comma-separated. Each must also belong to an IONOS CLOUD DNS zone you manage. Optional")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
