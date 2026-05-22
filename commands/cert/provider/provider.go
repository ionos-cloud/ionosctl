package provider

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cert "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			Use:              "provider",
			Aliases:          []string{"providers"},
			Short:            "Certificate Provider used to renew certificates before their expiry.",
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
	certClient := cert.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
	providers, _, err := certClient.ProviderApi.ProvidersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	t := table.New(allCols, table.WithPrefix("items"))
	if err := t.Extract(providers); err != nil {
		return nil
	}
	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").AddInfo("Email").AddInfo("Server").AddInfo("State").ToString()
}
