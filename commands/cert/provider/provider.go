package provider

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
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Email", JSONPath: "properties.email", Default: true},
	{Name: "Server", JSONPath: "properties.server", Default: true},
	{Name: "KeyId", JSONPath: "properties.externalAccountBinding.keyId", Default: true},
	{Name: "KeySecret", JSONPath: "properties.externalAccountBinding.keySecret", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ProviderCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "provider",
			Aliases: []string{"providers"},
			Short:   "Register ACME certificate authorities that issue and renew auto-certificates",
			Long: `A provider is an ACME (Automatic Certificate Management Environment) certificate authority registered in your account, for example Let's Encrypt. It records the CA's name, a contact email, the ACME directory URL, and optional external account binding (EAB) credentials that link your IONOS account to a pre-registered account at the CA.

Providers do not issue certificates on their own. Instead, an ` + "`autocertificate`" + ` references a provider by ID; the Certificate Manager then uses that provider to issue the certificate and to auto-renew it (~30 days before expiry) via ACME DNS-01 challenges against your IONOS CLOUD DNS zone.

Typical flow: create a provider, then create an autocertificate that points to it.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)
	cmd.AddCommand(ProviderPostCmd())
	cmd.AddCommand(ProviderDeleteCmd())
	cmd.AddCommand(ProviderFindByIdCmd())
	cmd.AddCommand(ProviderListCmd())
	cmd.AddCommand(ProviderPutCmd())

	return cmd
}

func ProviderIDs() []string {
	providers, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	t := table.New(allCols, table.WithPrefix("items"))
	if err := t.Extract(providers); err != nil {
		return nil
	}
	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").AddInfo("Email").AddInfo("Server").AddInfo("State").ToString()
}
