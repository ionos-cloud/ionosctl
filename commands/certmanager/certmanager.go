package certmanager

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/certmanager/certificate"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "certmanager",
			Aliases:          []string{"cert", "certs", "certificate-manager", "certificates", "certificate"},
			Short:            "The sub-commands of the 'certmanager' resource help you manage SSL certificates.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(CertGetApiVersionCmd())
	cmd.AddCommand(certificate.CertCmd())

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

	return core.WithRegionalFlags(cmd, constants.DNSApiRegionalURL, constants.DNSLocations)
}
