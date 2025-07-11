package certificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultCertificateCols = []string{"CertId", "DisplayName"}
)

func CertCmd() *core.Command {
	certCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "certificate",
			Aliases:          []string{"cert", "certificates", "certs"},
			Short:            "Certificate Manager Operations",
			Long:             "The sub-commands of `ionosctl certificate-manager` allows you to manage the SSL Certificates under your account.",
			TraverseChildren: true,
		},
	}

	certCmd.AddCommand(CertGetCmd())
	certCmd.AddCommand(CertCreateCmd())
	certCmd.AddCommand(CertListCmd())
	certCmd.AddCommand(CertDeleteCmd())
	certCmd.AddCommand(CertUpdateCmd())

	return certCmd
}

func printProperties(value cert.Certificate, flag bool) string {
	if !flag {
		return value.Certificate
	}
	return value.CertificateChain
}

func CertificatesIds() []string {
	ls, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(ls.Items, func(dto cert.CertificateRead) string {
		return dto.GetId()
	})
}
