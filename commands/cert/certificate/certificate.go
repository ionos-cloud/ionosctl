package certificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "CertId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.name", Default: true},
	{Name: "Expired", JSONPath: "metadata.expired", Default: true},
	{Name: "NotAfter", JSONPath: "metadata.notAfter", Default: true},
	{Name: "NotBefore", JSONPath: "metadata.notBefore", Default: true},
	{Name: "SerialNumber", JSONPath: "metadata.serialNumber"},
	{Name: "SubjectAlternativeNames", JSONPath: "metadata.subjectAlternativeNames"},
	{Name: "Certificate", JSONPath: "properties.certificate"},
	{Name: "CertificateChain", JSONPath: "properties.certificateChain"},
}

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

	certCmd.AddColsFlag(allCols)

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
	t := table.New(allCols, table.WithPrefix("items"))
	if err := t.Extract(ls); err != nil {
		return nil
	}
	return completions.NewCompleter(t.Rows(), "CertId").AddInfo("DisplayName").ToString()

}
