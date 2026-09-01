package autocertificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Provider", JSONPath: "properties.provider", Default: true},
	{Name: "CommonName", JSONPath: "properties.commonName", Default: true},
	{Name: "KeyAlgorithm", JSONPath: "properties.keyAlgorithm", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "AlternativeNames", JSONPath: "properties.subjectAlternativeNames", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func AutocertificateCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "autocertificate",
			Aliases: []string{"autocert", "autocerts", "auto", "autocertificates"},
			Short:   "Issue and auto-renew certificates via an ACME provider",
			Long: `An auto-certificate is a TLS/SSL certificate that IONOS issues and renews for you automatically, instead of one you upload yourself. It is tied to a provider (an ACME certificate authority, see ` + "`ionosctl certmanager provider`" + `) and to one or more domain names.

To issue the certificate, IONOS proves domain ownership using ACME DNS-01 challenges, which requires the certificate's common name (and any additional names) to belong to a DNS zone hosted in IONOS CLOUD DNS - IONOS creates the required TXT records for you. Renewal is automatic, roughly 30 days before expiry, with a grace period during which the old and new certificates temporarily coexist.

Prerequisites: a provider must already exist, and the domain(s) must be served by an IONOS CLOUD DNS zone you manage.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)
	cmd.AddCommand(AutocertificatePostCmd())
	cmd.AddCommand(AutocertificateDeleteCmd())
	cmd.AddCommand(AutocertificateFindByIdCmd())
	cmd.AddCommand(AutocertificateListCmd())
	cmd.AddCommand(AutocertificatePutCmd())

	return cmd
}

func AutocertificateIDs() []string {
	autocertificate, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	t := table.New(allCols, table.WithPrefix("items"))
	if err := t.Extract(autocertificate); err != nil {
		return nil
	}
	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").AddInfo("CommonName").AddInfo("KeyAlgorithm").AddInfo("Provider").AddInfo("AlternativeNames").AddInfo("State").ToString()
}
