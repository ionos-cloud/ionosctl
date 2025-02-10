package certmanager

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/services/certmanager/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/cobra"
)

var (
	defaultCertificateCols = []string{"CertId", "DisplayName"}
)

func CertCmd() *core.Command {
	certCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "certificate-manager",
			Short:            "Certificate Manager Operations",
			Long:             "The sub-commands of `ionosctl certificate-manager` allows you to manage the SSL Certificates under your account.",
			TraverseChildren: true,
		},
	}

	certCmd.AddCommand(CertGetCmd())
	certCmd.AddCommand(CertCreateCmd())
	certCmd.AddCommand(CertListCmd())
	certCmd.AddCommand(CertDeleteCmd())
	certCmd.AddCommand(CertGetApiVersionCmd())
	certCmd.AddCommand(CertUpdateCmd())

	return certCmd
}

func printProperties(value *cert.CertificateDto, flag bool) string {
	if !flag {
		return *value.Properties.Certificate
	}
	return *value.Properties.CertificateChain
}

func CertificatesIds() []string {
	client, _ := client2.Get()
	svc := resources.NewCertsService(client, context.Background())
	certs, _, _ := svc.List()
	return functional.Map(*certs.GetItems(), func(dto cert.CertificateDto) string {
		return *dto.GetId()
	})
}
