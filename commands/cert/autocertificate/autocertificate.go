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
			Use:              "autocertificate",
			Aliases:          []string{"autocert", "autocerts", "auto", "autocertificates"},
			Short:            "Auto certificates create new certificates based on a certificate provider",
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
