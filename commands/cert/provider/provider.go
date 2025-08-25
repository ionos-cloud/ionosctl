package provider

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Email", "Server", "KeyId", "KeySecret", "State"}
)

func ProviderCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "provider",
			Aliases:          []string{"p"},
			Short:            "Certificate Provider used to renew certificates before their expiry.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
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
	providersConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.CertManagerProvider, providers)
	if err != nil {
		return nil
	}
	return completions.NewCompleter(providersConverted, "Id").AddInfo("Name").AddInfo("Email").AddInfo("Server").AddInfo("State").ToString()
}
