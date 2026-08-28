package cert

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/cert/autocertificate"
	"github.com/ionos-cloud/ionosctl/v6/commands/cert/certificate"
	"github.com/ionos-cloud/ionosctl/v6/commands/cert/provider"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "certmanager",
			Aliases: []string{"cert", "certs", "certificate-manager", "certificates", "certificate"},
			Short:   "Manage TLS/SSL certificates for use by other IONOS products",
			Long: `The Certificate Manager stores TLS/SSL certificates so other IONOS products (Application Load Balancer, CDN, etc.) can serve HTTPS traffic under your domains. All certificates and private keys are handled in PEM (Privacy Enhanced Mail) format.

Certificates you upload and manage yourself (body + private key + chain) have no automatic renewal, so you replace them before expiry. Auto-certificates are instead issued and auto-renewed for you by an ACME provider (e.g. Let's Encrypt) you register, using DNS-01 challenges; this requires the domain's DNS zone to be hosted in IONOS CLOUD DNS.

Uploaded certificates need no DNS setup and work with ALB and CDN. Auto-certificates require IONOS CLOUD DNS and are auto-renewed.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(CertGetApiVersionCmd())
	cmd.AddCommand(certificate.CertCmd())
	cmd.AddCommand(provider.ProviderCommand())
	cmd.AddCommand(autocertificate.AutocertificateCommand())

	// Until v6.8.0, the following commands were included in the Root command.
	// We will still include them to avoid breaking changes, but they will be hidden
	// i.e. they will continue working but not show up in the help text anymore

	asHidden := func(cmd *core.Command) *core.Command {
		cmd.Command.Hidden = true
		return cmd
	}

	cmd.AddCommand(asHidden(certificate.CertCreateCmd()))
	cmd.AddCommand(asHidden(certificate.CertDeleteCmd()))
	cmd.AddCommand(asHidden(certificate.CertGetCmd()))
	cmd.AddCommand(asHidden(certificate.CertListCmd()))
	cmd.AddCommand(asHidden(certificate.CertUpdateCmd()))

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Cert}, constants.CertApiRegionalURL, constants.CertLocations)
}
